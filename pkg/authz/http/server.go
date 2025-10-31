package http

import (
	"context"
	"net/http"

	"github.com/google/cel-go/cel"
	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	kcel "github.com/kyverno/kyverno-envoy-plugin/pkg/cel"
	httpcel "github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/http"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/server"
	"k8s.io/client-go/dynamic"
)

type Config struct {
	NestedRequest    bool
	InputExpression  string
	OutputExpression string
}

func NewServer(addr string, dyn dynamic.Interface, p engine.HTTPSource, config Config) server.ServerFunc {
	// httpReq, err := httpcel.NewRequest(r)
	// ret, _, _ := prog.Eval(map[string]any{
	// 	"object": &httpReq,
	// })
	// fmt.Println(ret)
	return func(ctx context.Context) error {
		base, err := kcel.NewEnv(v1alpha1.EvaluationModeHTTP)
		if err != nil {
			return err
		}
		if config.InputExpression == "" {
			config.InputExpression = `
http.CheckRequest{
	method: object.Header("x-original-method")[0],
	header: object.header,
	host: url(object.Header("x-original-url")[0]).getHostname(),
	scheme: url(object.Header("x-original-url")[0]).getScheme(),
	path: url(object.Header("x-original-url")[0]).getEscapedPath(),
	query: url(object.Header("x-original-url")[0]).getQuery(),
	fragment: "todo",
}
`
		}
		inputEnv, err := base.Extend(cel.Variable("object", httpcel.RequestType))
		if err != nil {
			return err
		}
		inputAst, issues := inputEnv.Compile(config.InputExpression)
		if err := issues.Err(); err != nil {
			return err
		}
		inputProgram, err := inputEnv.Program(inputAst)
		if err != nil {
			return err
		}
		if config.OutputExpression == "" {
			config.OutputExpression = `
object.reason == ""
	? writer.WriteHeader(200)
	: writer.WriteHeader(401)
`
		}
		outputEnv, err := base.Extend(
			cel.Variable("object", httpcel.ResponseType),
			cel.Variable("writer", httpcel.ResponseWriterType),
		)
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
		// register health check
		a := &authorizer{
			provider:      p,
			dyn:           dyn,
			nestedRequest: config.NestedRequest,
			inputProgram:  inputProgram,
			outputProgram: outputProgram,
		}
		mux.Handle("POST /", a)
		// create server
		s := &http.Server{
			Addr:    addr,
			Handler: mux,
		}
		// run server
		return server.RunHttp(ctx, s, "", "")
	}
}
