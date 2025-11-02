package http

import (
	"bufio"
	"context"
	"fmt"
	"net/http"

	"github.com/google/cel-go/cel"
	httpcel "github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/authz/http"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
)

type authorizer struct {
	provider      engine.HTTPSource
	dyn           dynamic.Interface
	inputProgram  cel.Program
	outputProgram cel.Program
	nestedRequest bool
}

func (a *authorizer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctrl.LoggerFrom(r.Context()).Info("received request", "from", r.RemoteAddr)
	if a.nestedRequest {
		reader := bufio.NewReader(r.Body)
		req, err := http.ReadRequest(reader)
		if err != nil {
			writeErrResp(w, err)
			return
		}
		r = req
	}
	pols, err := a.provider.Load(context.Background())
	if err != nil {
		writeErrResp(w, err)
		return
	}
	httpReq, err := httpcel.NewRequest(r)
	if err != nil {
		writeErrResp(w, err)
		return
	}
	if a.inputProgram != nil {
		out, _, err := a.inputProgram.Eval(map[string]any{
			"object": &httpReq,
		})
		if err != nil {
			writeErrResp(w, err)
			return
		}
		if out.Value() != nil {
			out, ok := out.Value().(*httpcel.CheckRequest)
			if ok && out != nil {
				httpReq = *out
			}
		}
	}
	var result *httpcel.CheckResponse
	for _, pol := range pols {
		resp, err := pol.Evaluate(context.Background(), a.dyn, &httpReq)
		if err != nil {
			writeErrResp(w, err)
			return
		}
		// write the first valid policy response and exit
		if resp != nil {
			result = resp
			break
		}
	}
	if result == nil {
		result = &httpcel.CheckResponse{}
	}
	_, _, err = a.outputProgram.Eval(map[string]any{
		"object": result,
	})
	if err != nil {
		writeErrResp(w, err)
		return
	}
}

func writeErrResp(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, err.Error()) //nolint:errcheck
}
