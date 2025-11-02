package http

import (
	"io"
	"net/http"

	"github.com/google/cel-go/common/types"
)

var (
	RequestType           = types.NewObjectType("http.CheckRequest")
	RequestAttributesType = types.NewObjectType("http.CheckRequestAttributes")
	ResponseType          = types.NewObjectType("http.CheckResponse")
	ResponseOkType        = types.NewObjectType("http.CheckResponseOk")
	ResponseDeniedType    = types.NewObjectType("http.CheckResponseDenied")
)

type (
	header = map[string][]string
	query  = map[string][]string
)

type CheckRequestAttributes struct {
	Method        string `cel:"method"`
	Header        header `cel:"header"`
	Host          string `cel:"host"`
	Protocol      string `cel:"protocol"`
	ContentLength int64  `cel:"contentLength"`
	Body          []byte `cel:"body"`
	Scheme        string `cel:"scheme"`
	Path          string `cel:"path"`
	Query         query  `cel:"query"`
	Fragment      string `cel:"fragment"`
}

type CheckRequest struct {
	Attributes CheckRequestAttributes `cel:"attributes"`
}

type CheckResponseOk struct{}

type CheckResponseDenied struct {
	Reason string `cel:"reason"`
}

type CheckResponse struct {
	Ok     *CheckResponseOk     `cel:"ok"`
	Denied *CheckResponseDenied `cel:"denied"`
}

func NewRequest(r *http.Request) (CheckRequest, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return CheckRequest{}, err
	}
	return CheckRequest{
		Attributes: CheckRequestAttributes{
			Method:        r.Method,
			Header:        r.Header,
			Path:          r.URL.Path,
			Host:          r.Host,
			Protocol:      r.Proto,
			Body:          bodyBytes,
			Query:         r.URL.Query(),
			ContentLength: int64(len(bodyBytes)),
			Fragment:      r.URL.Fragment,
			Scheme:        r.URL.Scheme,
		},
	}, nil
}
