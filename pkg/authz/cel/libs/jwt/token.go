package jwt

import (
	"github.com/google/cel-go/common/types"
	"google.golang.org/protobuf/types/known/structpb"
)

var TokenType = types.NewObjectType("jwt.Token")

type Token struct {
	Header *structpb.Struct
	Claims *structpb.Struct
	Valid  bool
}
