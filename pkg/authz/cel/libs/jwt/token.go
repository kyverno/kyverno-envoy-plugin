package jwt

// import (
// 	"fmt"
// 	"reflect"

// 	"github.com/golang-jwt/jwt"
// 	"github.com/google/cel-go/common/types"
// 	"github.com/google/cel-go/common/types/ref"
// )

// var TokenType = types.NewObjectType("kyverno.jwt.Token")

// // Token provdes a CEL representation of a jwt.Token
// type Token struct {
// 	*jwt.Token
// }

// func (d Token) ConvertToNative(typeDesc reflect.Type) (interface{}, error) {
// 	if reflect.TypeOf(d.Token).AssignableTo(typeDesc) {
// 		return d.Token, nil
// 	}
// 	return nil, fmt.Errorf("type conversion error from 'Token' to '%v'", typeDesc)
// }

// func (d Token) ConvertToType(typeVal ref.Type) ref.Val {
// 	switch typeVal {
// 	case TokenType:
// 		return d
// 	case types.TypeType:
// 		return TokenType
// 	default:
// 		return types.NewErr("type conversion error from '%s' to '%s'", TokenType, typeVal)
// 	}
// }

// func (d Token) Equal(other ref.Val) ref.Val {
// 	otherToken, ok := other.(Token)
// 	if !ok {
// 		return types.MaybeNoSuchOverloadErr(other)
// 	}
// 	return types.Bool(d.Token.Raw == otherToken.Token.Raw)
// }

// func (d Token) Type() ref.Type {
// 	return TokenType
// }

// func (d Token) Value() any {
// 	return d.Token
// }
