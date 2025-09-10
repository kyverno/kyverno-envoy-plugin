package http

import (
	"io"
	"net/http"

	"github.com/google/cel-go/common/types"
)

var (
	RequestType  = types.NewObjectType("http.Request")
	KVType       = types.NewObjectType("http.KV")
	ResponseType = types.NewObjectType("http.Response")
)

type KV struct {
	inner map[string][]string `cel:"inner"`
}

func (k *KV) GetInnerMap() map[string][]string {
	return k.inner
}

type Request struct {
	Method   string `cel:"method"`
	Headers  *KV    `cel:"headers"`
	Path     string `cel:"path"`
	Host     string `cel:"host"`
	Scheme   string `cel:"scheme"`
	Query    *KV    `cel:"queryParams"`
	Fragment string `cel:"fragment"`
	Size     int64  `cel:"size"`
	Protocol string `cel:"protocol"`
	Body     string `cel:"body"`
	RawBody  []byte `cel:"rawBody"`
}

type Response struct {
	Status  int    `cel:"status"`
	Headers *KV    `cel:"headers"`
	Body    string `cel:"body"`
}

func NewRequest(r *http.Request) (Request, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return Request{}, err
	}
	return Request{
		Method:   r.Method,
		Headers:  &KV{inner: r.Header},
		Path:     r.URL.Path,
		Host:     r.URL.Host,
		Protocol: r.Proto,
		RawBody:  bodyBytes,
		Body:     string(bodyBytes),
		Query:    &KV{inner: r.URL.Query()},
		Size:     int64(len(bodyBytes)),
		Fragment: r.URL.Fragment,
		Scheme:   r.URL.Scheme,
	}, nil
}
