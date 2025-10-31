package http

import (
	"io"
	"net/http"

	"github.com/google/cel-go/common/types"
)

var (
	RequestType  = types.NewObjectType("http.CheckRequest")
	ResponseType = types.NewObjectType("http.CheckResponse")
	// KVType       = types.NewObjectType("http.KV")
)

type Header = map[string][]string
type Query = map[string][]string

type CheckRequest struct {
	// from request
	Method        string              `cel:"method"`
	Header        map[string][]string `cel:"header"`
	Host          string              `cel:"host"`
	Protocol      string              `cel:"protocol"`
	ContentLength int64               `cel:"contentLength"`
	Body          string              `cel:"body"`
	RawBody       []byte              `cel:"rawBody"`
	// from url
	Scheme   string              `cel:"scheme"`
	Path     string              `cel:"path"`
	Query    map[string][]string `cel:"query"`
	Fragment string              `cel:"fragment"`
}

type CheckResponse struct {
	Status int                 `cel:"status"`
	Header map[string][]string `cel:"header"`
	Body   string              `cel:"body"`
}

func NewRequest(r *http.Request) (CheckRequest, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return CheckRequest{}, err
	}
	return CheckRequest{
		Method:        r.Method,
		Header:        r.Header,
		Path:          r.URL.Path,
		Host:          r.Host,
		Protocol:      r.Proto,
		RawBody:       bodyBytes,
		Body:          string(bodyBytes),
		Query:         r.URL.Query(),
		ContentLength: int64(len(bodyBytes)),
		Fragment:      r.URL.Fragment,
		Scheme:        r.URL.Scheme,
	}, nil
}
