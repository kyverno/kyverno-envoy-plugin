package envoy

import (
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typesv3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/authz/cel/utils"
	status "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/structpb"
)

type impl struct {
	types.Adapter
}

func (c *impl) allowed() ref.Val {
	r := &authv3.OkHttpResponse{}
	return c.NativeToValue(r)
}

func (c *impl) ok_with_header(ok ref.Val, header ref.Val) ref.Val {
	if ok, err := utils.ConvertToNative[*authv3.OkHttpResponse](ok); err != nil {
		return types.WrapErr(err)
	} else if header, err := utils.ConvertToNative[*corev3.HeaderValueOption](header); err != nil {
		return types.WrapErr(err)
	} else {
		ok.Headers = append(ok.Headers, header)
		return c.NativeToValue(ok)
	}
}

func (c *impl) ok_without_header(ok ref.Val, header ref.Val) ref.Val {
	if ok, err := utils.ConvertToNative[*authv3.OkHttpResponse](ok); err != nil {
		return types.WrapErr(err)
	} else if header, err := utils.ConvertToNative[string](header); err != nil {
		return types.WrapErr(err)
	} else {
		ok.HeadersToRemove = append(ok.HeadersToRemove, header)
		return c.NativeToValue(ok)
	}
}

func (c *impl) ok_with_response_header(ok ref.Val, header ref.Val) ref.Val {
	if ok, err := utils.ConvertToNative[*authv3.OkHttpResponse](ok); err != nil {
		return types.WrapErr(err)
	} else if header, err := utils.ConvertToNative[*corev3.HeaderValueOption](header); err != nil {
		return types.WrapErr(err)
	} else {
		ok.ResponseHeadersToAdd = append(ok.ResponseHeadersToAdd, header)
		return c.NativeToValue(ok)
	}
}

func (c *impl) ok_with_query_param(ok ref.Val, param ref.Val) ref.Val {
	if ok, err := utils.ConvertToNative[*authv3.OkHttpResponse](ok); err != nil {
		return types.WrapErr(err)
	} else if param, err := utils.ConvertToNative[*corev3.QueryParameter](param); err != nil {
		return types.WrapErr(err)
	} else {
		ok.QueryParametersToSet = append(ok.QueryParametersToSet, param)
		return c.NativeToValue(ok)
	}
}

func (c *impl) ok_without_query_param(ok ref.Val, param ref.Val) ref.Val {
	if ok, err := utils.ConvertToNative[*authv3.OkHttpResponse](ok); err != nil {
		return types.WrapErr(err)
	} else if param, err := utils.ConvertToNative[string](param); err != nil {
		return types.WrapErr(err)
	} else {
		ok.QueryParametersToRemove = append(ok.QueryParametersToRemove, param)
		return c.NativeToValue(ok)
	}
}

func (c *impl) denied(code ref.Val) ref.Val {
	if code, err := utils.ConvertToNative[typesv3.StatusCode](code); err != nil {
		return types.WrapErr(err)
	} else {
		return c.NativeToValue(&authv3.DeniedHttpResponse{Status: &typesv3.HttpStatus{Code: code}})
	}
}

func (c *impl) denied_with_body(denied ref.Val, body ref.Val) ref.Val {
	if denied, err := utils.ConvertToNative[*authv3.DeniedHttpResponse](denied); err != nil {
		return types.WrapErr(err)
	} else if body, err := utils.ConvertToNative[string](body); err != nil {
		return types.WrapErr(err)
	} else {
		denied.Body = body
		return c.NativeToValue(denied)
	}
}

func (c *impl) denied_with_header(denied ref.Val, header ref.Val) ref.Val {
	if denied, err := utils.ConvertToNative[*authv3.DeniedHttpResponse](denied); err != nil {
		return types.WrapErr(err)
	} else if header, err := utils.ConvertToNative[*corev3.HeaderValueOption](header); err != nil {
		return types.WrapErr(err)
	} else {
		denied.Headers = append(denied.Headers, header)
		return c.NativeToValue(denied)
	}
}

func (c *impl) header_key_value(key ref.Val, value ref.Val) ref.Val {
	if key, err := utils.ConvertToNative[string](key); err != nil {
		return types.WrapErr(err)
	} else if value, err := utils.ConvertToNative[string](value); err != nil {
		return types.WrapErr(err)
	} else {
		return c.NativeToValue(&corev3.HeaderValueOption{Header: &corev3.HeaderValue{Key: key, Value: value}})
	}
}

func (c *impl) header_keep_empty_value(header ref.Val) ref.Val {
	return c.header_keep_empty_value_bool(header, types.True)
}

func (c *impl) header_keep_empty_value_bool(header ref.Val, flag ref.Val) ref.Val {
	if header, err := utils.ConvertToNative[*corev3.HeaderValueOption](header); err != nil {
		return types.WrapErr(err)
	} else if flag, err := utils.ConvertToNative[bool](flag); err != nil {
		return types.WrapErr(err)
	} else {
		header.KeepEmptyValue = flag
		return c.NativeToValue(header)
	}
}

func (c *impl) response_code(code ref.Val) ref.Val {
	if code, err := utils.ConvertToNative[codes.Code](code); err != nil {
		return types.WrapErr(err)
	} else {
		return c.NativeToValue(&authv3.CheckResponse{
			Status: &status.Status{Code: int32(code)},
		})
	}
}

func (c *impl) response_ok(ok ref.Val) ref.Val {
	if ok, err := utils.ConvertToNative[*authv3.OkHttpResponse](ok); err != nil {
		return types.WrapErr(err)
	} else {
		return c.NativeToValue(&authv3.CheckResponse{
			Status:       &status.Status{Code: int32(codes.OK)},
			HttpResponse: &authv3.CheckResponse_OkResponse{OkResponse: ok},
		})
	}
}

func (c *impl) response_denied(denied ref.Val) ref.Val {
	if denied, err := utils.ConvertToNative[*authv3.DeniedHttpResponse](denied); err != nil {
		return types.WrapErr(err)
	} else {
		return c.NativeToValue(&authv3.CheckResponse{
			Status:       &status.Status{Code: int32(codes.PermissionDenied)},
			HttpResponse: &authv3.CheckResponse_DeniedResponse{DeniedResponse: denied},
		})
	}
}

func (c *impl) response_with_message(response ref.Val, message ref.Val) ref.Val {
	if response, err := utils.ConvertToNative[*authv3.CheckResponse](response); err != nil {
		return types.WrapErr(err)
	} else if message, err := utils.ConvertToNative[string](message); err != nil {
		return types.WrapErr(err)
	} else {
		response.Status.Message = message
		return c.NativeToValue(response)
	}
}

func (c *impl) response_with_metadata(response ref.Val, metadata ref.Val) ref.Val {
	if response, err := utils.ConvertToNative[*authv3.CheckResponse](response); err != nil {
		return types.WrapErr(err)
	} else if metadata, err := utils.ConvertToNative[*structpb.Struct](metadata); err != nil {
		return types.WrapErr(err)
	} else {
		response.DynamicMetadata = metadata
		return c.NativeToValue(response)
	}
}
