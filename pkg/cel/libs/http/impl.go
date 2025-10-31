package http

import (
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/cel/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

func (c *impl) get_header_value(allHeaders ref.Val, header ref.Val) ref.Val {
	if kv, err := utils.ConvertToNative[KV](allHeaders); err != nil {
		return types.WrapErr(err)
	} else if header, err := utils.ConvertToNative[string](header); err != nil {
		return types.WrapErr(err)
	} else {
		caser := cases.Title(language.Und) // turn all instances of a header to match a single case
		v, exists := kv[caser.String(header)]
		if exists {
			return c.NativeToValue(v[0])
		}
		v, exists = kv[header] // try to get the header directly
		if exists {
			return c.NativeToValue(v[0])
		}
		return c.NativeToValue("")
	}
}

func (c *impl) get_header_all(allHeaders ref.Val, header ref.Val) ref.Val {
	if kv, err := utils.ConvertToNative[KV](allHeaders); err != nil {
		return types.WrapErr(err)
	} else if header, err := utils.ConvertToNative[string](header); err != nil {
		return types.WrapErr(err)
	} else {
		caser := cases.Title(language.Und)
		v, exists := kv[caser.String(header)]
		if exists {
			return c.NativeToValue(v)
		}
		v, exists = kv[header]
		if exists {
			return c.NativeToValue(v)
		}
		return c.NativeToValue([]string{})
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
		if r.Headers == nil {
			r.Headers = KV{}
		}
		r.Headers[k] = append(r.Headers[k], v)
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
