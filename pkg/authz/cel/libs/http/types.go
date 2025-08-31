package http

import (
	"net/http"

	"github.com/google/cel-go/common/types"
)

var (
	RequestType  = types.NewObjectType("http.Request")
	KVType       = types.NewObjectType("http.KV")
	ResponseType = types.NewObjectType("http.response")
)

type responseProvider struct {
	types.Provider
}

type KV struct {
	kv map[string][]string
}

func NewResponseProvider(p types.Provider) *responseProvider {
	return &responseProvider{
		Provider: p,
	}
}

func (p *responseProvider) FindStructType(typeName string) (*types.Type, bool) {
	if typeName == "http.response" {
		return p.Provider.FindStructType("http.Response")
	}
	return p.Provider.FindStructType(typeName)
}

func (p *responseProvider) FindStructFieldType(typeName, fieldName string) (*types.FieldType, bool) {
	if typeName == "http.response" {
		return p.Provider.FindStructFieldType("http.Response", fieldName)
	}
	return p.Provider.FindStructFieldType(typeName, fieldName)
}

type Request struct {
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

type Response struct {
	Status  int    `cel:"status"`
	Headers KV     `cel:"headers"`
	Body    string `cel:"body"`
}

// flow -> compiled policies received from earlier -> request received ->

func ToCELRequest(r http.Request) Request {
	return Request{
		Method:  r.Method,
		Headers: KV{kv: r.Header},
		Path:    r.URL.Path,
		Host:    r.URL.Host,
	}
}

func ToNativeResponse(r Response) http.Response {
	return http.Response{}
}
