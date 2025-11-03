package http

import (
	"context"
	"net/http"

	"github.com/google/cel-go/cel"
	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	kcel "github.com/kyverno/kyverno-envoy-plugin/pkg/cel"
	httpcel "github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/authz/http"
	httpserver "github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/httpserver"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/server"
	"k8s.io/client-go/dynamic"
)

func NewServer(config Config, p engine.HTTPSource, dyn dynamic.Interface) server.ServerFunc {
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
		// register service
		a := &authorizer{
			provider:      p,
			dyn:           dyn,
			inputProgram:  inputProgram,
			outputProgram: outputProgram,
			nestedRequest: config.NestedRequest,
		}
		mux.Handle("POST /", a)
		// create server
		s := &http.Server{
			Addr:    config.Address,
			Handler: mux,
		}
		// run server
		return server.RunHttp(ctx, s, "", "")
	}
}
