package httpserver

import (
	"github.com/google/cel-go/common/types"
)

var (
	ResponseType = types.NewObjectType("httpserver.HttpResponse")
)

type HttpResponse struct {
	Status int                 `cel:"status"`
	Header map[string][]string `cel:"header"`
	Body   []byte              `cel:"body"`
}
