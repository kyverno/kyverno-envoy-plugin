package httpauth

import (
	"bufio"
	"context"
	"fmt"
	"net/http"

	httpcel "github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel/libs/http"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
)

type Authorizer struct {
	provider engine.Provider
}

func (a *Authorizer) NewHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		reader := bufio.NewReader(r.Body)
		req, err := http.ReadRequest(reader)
		if err != nil {
			writeErrResp(w, err)
			return
		}
		pols, err := a.provider.CompiledPolicies(context.Background())
		if err != nil {
			writeErrResp(w, err)
			return
		}
		ruleFuncs := []engine.RequestFunc{}
		for _, pol := range pols {
			ruleFuncs = append(ruleFuncs, pol.ForHTTP(req))
		}
		for _, r := range ruleFuncs {
			resp, err := r()
			if err != nil {
				writeErrResp(w, err)
				return
			}
			// write the first valid policy response and exit
			if resp != nil {
				writeResponse(w, resp)
				break
			}
		}
	}
}

func writeErrResp(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, err.Error())
}

func writeResponse(w http.ResponseWriter, resp *httpcel.Response) {
	for k, v := range resp.Headers.GetInnerMap() {
		for _, val := range v {
			w.Header().Set(k, val)
		}
	}
	w.WriteHeader(resp.Status)
	fmt.Fprint(w, resp.Body)
}
