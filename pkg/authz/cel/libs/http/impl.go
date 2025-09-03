package http

import (
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type impl struct {
	types.Adapter
}

func (c *impl) response() ref.Val {
	r := &Response{}
	return c.NativeToValue(r)
}

func (c *impl) get_header_value(req ref.Val, header ref.Val) ref.Val {
	if kv, err := utils.ConvertToNative[KV](req); err != nil {
		return types.WrapErr(err)
	} else if header, err := utils.ConvertToNative[string](header); err != nil {
		return types.WrapErr(err)
	} else {
		caser := cases.Title(language.Und) // turn all instances of a header to match a single case
		v, exists := kv.inner[caser.String(header)]
		if !exists {
			return c.NativeToValue("")
		}
		return c.NativeToValue(v[0])
	}
}

func (c *impl) get_header_all(req ref.Val, header ref.Val) ref.Val {
	if kv, err := utils.ConvertToNative[KV](req); err != nil {
		return types.WrapErr(err)
	} else if header, err := utils.ConvertToNative[string](header); err != nil {
		return types.WrapErr(err)
	} else {
		v, exists := kv.inner[header]
		if !exists {
			return c.NativeToValue([]string{})
		}
		return c.NativeToValue(v)
	}
}

func (c *impl) with_status(r ref.Val, status ref.Val) ref.Val {
	if r, err := utils.ConvertToNative[*Response](r); err != nil {
		return types.WrapErr(err)
	} else if statusCode, err := utils.ConvertToNative[int](status); err != nil {
		return types.WrapErr(err)
	} else {
		r.Status = statusCode
		return c.NativeToValue(r)
	}
}

func (c *impl) with_header(args ...ref.Val) ref.Val {
	if r, err := utils.ConvertToNative[*Response](args[0]); err != nil {
		return types.WrapErr(err)
	} else if k, err := utils.ConvertToNative[string](args[1]); err != nil {
		return types.WrapErr(err)
	} else if v, err := utils.ConvertToNative[string](args[2]); err != nil {
		return types.WrapErr(err)
	} else {
		r.Headers.inner[k] = append(r.Headers.inner[k], v)
		return c.NativeToValue(r)
	}
}

func (c *impl) with_body(r ref.Val, b ref.Val) ref.Val {
	if r, err := utils.ConvertToNative[*Response](r); err != nil {
		return types.WrapErr(err)
	} else if b, err := utils.ConvertToNative[string](b); err != nil {
		return types.WrapErr(err)
	} else {
		r.Body = b
		return c.NativeToValue(r)
	}
}
