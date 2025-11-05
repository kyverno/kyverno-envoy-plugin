package http

import (
	"bufio"
	"fmt"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/cel-go/cel"
	httpcel "github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/authz/http"
	httpserver "github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/httpserver"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/cel/utils"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/metrics"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/extensions/policy"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
)

type authorizer struct {
	engine        core.Engine[dynamic.Interface, *httpcel.CheckRequest, policy.Evaluation[*httpcel.CheckResponse]]
	dyn           dynamic.Interface
	inputProgram  cel.Program
	outputProgram cel.Program
	nestedRequest bool
}

func (a *authorizer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	logger := ctrl.LoggerFrom(r.Context()).WithValues("from", r.RemoteAddr)
	logger.Info("received request")
	if a.nestedRequest {
		reader := bufio.NewReader(r.Body)
		req, err := http.ReadRequest(reader)
		if err != nil {
			writeErrResp(w, err)
			return
		}
		r = req
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
	response := a.engine.Handle(r.Context(), a.dyn, &httpReq)
	if response.Error != nil {
		metrics.RecordHTTPRequestError(r.Context(), httpReq, response.Error)
		writeErrResp(w, response.Error)
		return
	}
	result := response.Result
	if result == nil {
		result = &httpcel.CheckResponse{
			Ok: &httpcel.CheckResponseOk{},
		}
	}
	defer metrics.RecordHTTPRequest(r.Context(), start, httpReq, result)
	out, _, err := a.outputProgram.Eval(map[string]any{
		"object": result,
	})
	if err != nil {
		writeErrResp(w, err)
		return
	}
	if out, err := utils.ConvertToNative[httpserver.HttpResponse](out); err != nil {
		writeErrResp(w, err)
	} else {
		writeResponse(logger, w, out)
	}
}

func writeErrResp(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, err.Error()) //nolint:errcheck
}

func writeResponse(logger logr.Logger, w http.ResponseWriter, resp httpserver.HttpResponse) {
	if resp.Header != nil {
		for k, v := range resp.Header {
			for _, val := range v {
				w.Header().Set(k, val)
			}
		}
	}
	w.WriteHeader(resp.Status)
	if resp.Body != nil {
		_, err := w.Write(resp.Body)
		if err != nil {
			logger.Error(err, "failed to write body")
		}
	}
}
