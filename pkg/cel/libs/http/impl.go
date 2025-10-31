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

func (c *impl) allowed() ref.Val {
	r := &CheckResponse{}
	return c.NativeToValue(r)
}

func (c *impl) denied(reason ref.Val) ref.Val {
	if reason, err := utils.ConvertToNative[string](reason); err != nil {
		return types.WrapErr(err)
	} else {
		r := &CheckResponse{
			Reason: reason,
		}
		return c.NativeToValue(r)
	}
}

// func (c *impl) with_header(args ...ref.Val) ref.Val {
// 	if r, err := utils.ConvertToNative[*CheckResponse](args[0]); err != nil {
// 		return types.WrapErr(err)
// 	} else if k, err := utils.ConvertToNative[string](args[1]); err != nil {
// 		return types.WrapErr(err)
// 	} else if v, err := utils.ConvertToNative[string](args[2]); err != nil {
// 		return types.WrapErr(err)
// 	} else {
// 		if r.Header == nil {
// 			r.Header = Header{}
// 		}
// 		r.Header[k] = append(r.Header[k], v)
// 		return c.NativeToValue(r)
// 	}
// }

// func (c *impl) with_body(r ref.Val, b ref.Val) ref.Val {
// 	if r, err := utils.ConvertToNative[*CheckResponse](r); err != nil {
// 		return types.WrapErr(err)
// 	} else if b, err := utils.ConvertToNative[string](b); err != nil {
// 		return types.WrapErr(err)
// 	} else {
// 		r.Body = b
// 		return c.NativeToValue(r)
// 	}
// }

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

func (c *impl) write(lhs ref.Val, rhs ref.Val) ref.Val {
	if writer, err := utils.ConvertToNative[ResponseWriter](lhs); err != nil {
		return types.WrapErr(err)
	} else if bytes, err := utils.ConvertToNative[[]byte](rhs); err != nil {
		return types.WrapErr(err)
	} else if _, err := writer.Write(bytes); err != nil {
		return types.WrapErr(err)
	} else {
		return c.NativeToValue(writer)
	}
}

func (c *impl) write_header(lhs ref.Val, rhs ref.Val) ref.Val {
	if writer, err := utils.ConvertToNative[ResponseWriter](lhs); err != nil {
		return types.WrapErr(err)
	} else if status, err := utils.ConvertToNative[int](rhs); err != nil {
		return types.WrapErr(err)
	} else {
		writer.WriteHeader(status)
		return c.NativeToValue(writer)
	}
}
