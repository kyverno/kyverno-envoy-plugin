package httpauth

import (
	"bufio"
	"context"
	"fmt"
	"net/http"

	httpauth "github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/http"
	httpcel "github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/http"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/dynamic"
)

type authorizer struct {
	provider      engine.HTTPSource
	logger        *logrus.Logger
	dyn           dynamic.Interface
	nestedRequest bool
}

func NewAuthorizer(dyn dynamic.Interface, p engine.HTTPSource, nestedRequest bool, logger *logrus.Logger) *authorizer {
	return &authorizer{
		provider:      p,
		logger:        logger,
		dyn:           dyn,
		nestedRequest: nestedRequest,
	}
}

func (a *authorizer) NewHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		a.logger.Infof("received request from %s", r.RemoteAddr)
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
		httpReq, err := httpauth.NewRequest(r)
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
