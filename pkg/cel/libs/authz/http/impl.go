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
	r := CheckResponseOk{}
	return c.NativeToValue(r)
}

func (c *impl) denied(reason ref.Val) ref.Val {
	if reason, err := utils.ConvertToNative[string](reason); err != nil {
		return types.WrapErr(err)
	} else {
		r := CheckResponseDenied{
			Reason: reason,
		}
		return c.NativeToValue(r)
	}
}

func (c *impl) get_header(request ref.Val, key ref.Val) ref.Val {
	if request, err := utils.ConvertToNative[CheckRequestAttributes](request); err != nil {
		return types.WrapErr(err)
	} else if key, err := utils.ConvertToNative[string](key); err != nil {
		return types.WrapErr(err)
	} else {
		return c.NativeToValue(textproto.MIMEHeader(request.Header).Values(key))
	}
}

func (c *impl) get_queryparam(request ref.Val, key ref.Val) ref.Val {
	if request, err := utils.ConvertToNative[CheckRequestAttributes](request); err != nil {
		return types.WrapErr(err)
	} else if key, err := utils.ConvertToNative[string](key); err != nil {
		return types.WrapErr(err)
	} else {
		return c.NativeToValue(request.Query[key])
	}
}

func (c *impl) response_ok(response ref.Val) ref.Val {
	if response, err := utils.ConvertToNative[CheckResponseOk](response); err != nil {
		return types.WrapErr(err)
	} else {
		r := &CheckResponse{
			Ok: &response,
		}
		return c.NativeToValue(r)
	}
}

func (c *impl) response_denied(response ref.Val) ref.Val {
	if response, err := utils.ConvertToNative[CheckResponseDenied](response); err != nil {
		return types.WrapErr(err)
	} else {
		r := &CheckResponse{
			Denied: &response,
		}
		return c.NativeToValue(r)
	}
}
