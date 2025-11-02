package httpserver

import (
	"github.com/google/cel-go/common/types"
)

var (
	ResponseType = types.NewObjectType("httpserver.CheckResponse")
)

type CheckResponse struct {
	Status int                 `cel:"status"`
	Header map[string][]string `cel:"header"`
	Body   []byte              `cel:"body"`
}
