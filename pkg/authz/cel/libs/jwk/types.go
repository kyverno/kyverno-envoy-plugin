package jwk

import (
	"github.com/google/cel-go/common/types"
	"github.com/lestrrat-go/jwx/v3/jwk"
)

var SetType = types.NewOpaqueType("jwk.Set")

type Set struct {
	jwk.Set
}
