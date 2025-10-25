package httpauth

import (
	"bufio"
	"context"
	"fmt"
	"net/http"

	httpcel "github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/http"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
)

type authorizer struct {
	provider      engine.HTTPSource
	dyn           dynamic.Interface
	nestedRequest bool
}

func (a *authorizer) NewHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
		for _, pol := range pols {
			resp, err := pol.Evaluate(context.Background(), a.dyn, &httpReq)
			if err != nil {
				writeErrResp(w, err)
				return
			}
			// write the first valid policy response and exit
			if resp != nil {
				writeResponse(w, resp)
				return
			}
		}
	}
}

func writeErrResp(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, err.Error()) //nolint:errcheck
}

func writeResponse(w http.ResponseWriter, resp *httpcel.Response) {
	if resp.Headers != nil {
		for k, v := range resp.Headers.GetInnerMap() {
			for _, val := range v {
				w.Header().Set(k, val)
			}
		}
	}

	w.WriteHeader(resp.Status)
	fmt.Fprint(w, resp.Body) //nolint:errcheck
}
