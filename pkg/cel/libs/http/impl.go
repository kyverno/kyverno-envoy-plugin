package http

import (
	"net/textproto"

	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/cel/utils"
)

type impl struct {
	types.Adapter
}

func (c *impl) response(statusCode ref.Val) ref.Val {
	if status, err := utils.ConvertToNative[int](statusCode); err != nil {
		return types.WrapErr(err)
	} else {
		r := &CheckResponse{
			Status: status,
		}
		return c.NativeToValue(r)
	}
}

func (c *impl) with_header(args ...ref.Val) ref.Val {
	if r, err := utils.ConvertToNative[*CheckResponse](args[0]); err != nil {
		return types.WrapErr(err)
	} else if k, err := utils.ConvertToNative[string](args[1]); err != nil {
		return types.WrapErr(err)
	} else if v, err := utils.ConvertToNative[string](args[2]); err != nil {
		return types.WrapErr(err)
	} else {
		if r.Header == nil {
			r.Header = Header{}
		}
		r.Header[k] = append(r.Header[k], v)
		return c.NativeToValue(r)
	}
}

func (c *impl) with_body(r ref.Val, b ref.Val) ref.Val {
	if r, err := utils.ConvertToNative[*CheckResponse](r); err != nil {
		return types.WrapErr(err)
	} else if b, err := utils.ConvertToNative[string](b); err != nil {
		return types.WrapErr(err)
	} else {
		r.Body = b
		return c.NativeToValue(r)
	}
}

func (c *impl) get_header(lhs ref.Val, rhs ref.Val) ref.Val {
	if request, err := utils.ConvertToNative[*CheckRequest](lhs); err != nil {
		return types.WrapErr(err)
	} else if key, err := utils.ConvertToNative[string](rhs); err != nil {
		return types.WrapErr(err)
	} else {
		return c.NativeToValue(textproto.MIMEHeader(request.Header).Values(key))
	}
}

func (c *impl) get_queryparam(lhs ref.Val, rhs ref.Val) ref.Val {
	if request, err := utils.ConvertToNative[*CheckRequest](lhs); err != nil {
		return types.WrapErr(err)
	} else if key, err := utils.ConvertToNative[string](rhs); err != nil {
		return types.WrapErr(err)
	} else {
		return c.NativeToValue(request.Query[key])
	}
}
