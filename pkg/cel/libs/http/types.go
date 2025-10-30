package http

import (
	"io"
	"net/http"

	"github.com/google/cel-go/common/types"
)

var (
	RequestType  = types.NewObjectType("http.Req")
	KVType       = types.NewObjectType("http.KV")
	ResponseType = types.NewObjectType("http.Resp")
)

type KV map[string][]string

type Req struct {
	Method   string `cel:"method"`
	Headers  KV     `cel:"headers"`
	Path     string `cel:"path"`
	Host     string `cel:"host"`
	Scheme   string `cel:"scheme"`
	Query    KV     `cel:"queryParams"`
	Fragment string `cel:"fragment"`
	Size     int64  `cel:"size"`
	Protocol string `cel:"protocol"`
	Body     string `cel:"body"`
	RawBody  []byte `cel:"rawBody"`
}

type Resp struct {
	Status  int    `cel:"status"`
	Headers KV     `cel:"headers"`
	Body    string `cel:"body"`
}

func NewRequest(r *http.Request) (Req, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return Req{}, err
	}
	return Req{
		Method:   r.Method,
		Headers:  KV(r.Header),
		Path:     r.URL.Path,
		Host:     r.Host,
		Protocol: r.Proto,
		RawBody:  bodyBytes,
		Body:     string(bodyBytes),
		Query:    KV(r.URL.Query()),
		Size:     int64(len(bodyBytes)),
		Fragment: r.URL.Fragment,
		Scheme:   r.URL.Scheme,
	}, nil
}
