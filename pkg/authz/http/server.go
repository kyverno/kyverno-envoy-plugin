package http

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/google/cel-go/cel"
	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	kcel "github.com/kyverno/kyverno-envoy-plugin/pkg/cel"
	httpcel "github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/authz/http"
	httpserver "github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/httpserver"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/server"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core/dispatchers"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core/handlers"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core/resulters"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/extensions/policy"
	"k8s.io/client-go/dynamic"
)

func NewServer(config Config, source engine.HTTPSource, dyn dynamic.Interface) server.ServerFunc {
	return func(ctx context.Context) error {
		base, err := kcel.NewEnv(v1alpha1.EvaluationModeHTTP)
		if err != nil {
			return err
		}
		var inputProgram cel.Program
		if config.InputExpression != "" {
			// 			config.InputExpression = `
			// http.CheckRequest{
			// 	attributes: http.CheckRequestAttributes{
			// 		method: object.attributes.Header("x-original-method")[0],
			// 		header: object.attributes.header,
			// 		host: url(object.attributes.Header("x-original-url")[0]).getHostname(),
			// 		scheme: url(object.attributes.Header("x-original-url")[0]).getScheme(),
			// 		path: url(object.attributes.Header("x-original-url")[0]).getEscapedPath(),
			// 		query: url(object.attributes.Header("x-original-url")[0]).getQuery(),
			// 		body: object.attributes.body,
			// 		fragment: "todo",
			// 	}
			// }
			// `
			inputEnv, err := base.Extend(cel.Variable("object", httpcel.RequestType))
			if err != nil {
				return err
			}
			inputAst, issues := inputEnv.Compile(config.InputExpression)
			if err := issues.Err(); err != nil {
				return err
			}
			program, err := inputEnv.Program(inputAst)
			if err != nil {
				return err
			}
			inputProgram = program
		}
		if config.OutputExpression == "" {
			config.OutputExpression = `
has(object.ok)
	? httpserver.HttpResponse{ status: 200 }
	: httpserver.HttpResponse{ status: 403, body: bytes(object.denied.reason) }
`
		}
		outputEnv, err := base.Extend(
			cel.Variable("object", httpcel.ResponseType),
			httpserver.Lib(),
		)
		if err != nil {
			return err
		}
		outputAst, issues := outputEnv.Compile(config.OutputExpression)
		if err := issues.Err(); err != nil {
			return err
		}
		outputProgram, err := outputEnv.Program(outputAst)
		if err != nil {
			return err
		}
		// create mux
		mux := http.NewServeMux()
		// build the engine
		engine := core.NewEngine(
			source,
			handlers.Handler(
				dispatchers.Sequential(
					policy.EvaluatorFactory[engine.HTTPPolicy](),
					func(ctx context.Context, fc core.FactoryContext[engine.HTTPPolicy, dynamic.Interface, *httpcel.CheckRequest]) core.Breaker[engine.HTTPPolicy, *httpcel.CheckRequest, policy.Evaluation[*httpcel.CheckResponse]] {
						return core.MakeBreakerFunc(func(_ context.Context, _ engine.HTTPPolicy, _ *httpcel.CheckRequest, out policy.Evaluation[*httpcel.CheckResponse]) bool {
							return out.Result != nil
						})
					},
				),
				func(ctx context.Context, fc core.FactoryContext[engine.HTTPPolicy, dynamic.Interface, *httpcel.CheckRequest]) core.Resulter[engine.HTTPPolicy, *httpcel.CheckRequest, policy.Evaluation[*httpcel.CheckResponse], policy.Evaluation[*httpcel.CheckResponse]] {
					return resulters.NewFirst[engine.HTTPPolicy, *httpcel.CheckRequest](func(out policy.Evaluation[*httpcel.CheckResponse]) bool {
						return out.Result != nil || out.Error != nil
					})
				},
			),
		)
		// register service
		a := &authorizer{
			engine:        engine,
			dyn:           dyn,
			inputProgram:  inputProgram,
			outputProgram: outputProgram,
			nestedRequest: config.NestedRequest,
		}
		mux.Handle("POST /{$}", a)
		// create server
		s := &http.Server{
			Addr:    config.Address,
			Handler: mux,
		}
		// serve TLS if a certfile and a keyfile are provided
		if config.CertFile != "" && config.KeyFile != "" {
			s.TLSConfig = &tls.Config{
				MinVersion: tls.VersionTLS12,
				CipherSuites: []uint16{
					// AEADs w/ ECDHE
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
					tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				},
			}
		}
		// run server
		return server.RunHttp(ctx, s, config.CertFile, config.KeyFile)
	}
}
