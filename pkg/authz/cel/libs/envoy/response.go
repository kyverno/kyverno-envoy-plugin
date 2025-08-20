package envoy

import (
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	status "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type OkResponse struct {
	// Status “OK“ allows the request. Any other status indicates the request should be denied, and
	// for HTTP filter, if not overridden by :ref:`denied HTTP response status <envoy_v3_api_field_service.auth.v3.DeniedHttpResponse.status>`
	// Envoy sends “403 Forbidden“ HTTP status code by default.
	Status *status.Status `cel:"status"`
	// An message that contains HTTP response attributes. This message is
	// used when the authorization service needs to send custom responses to the
	// downstream client or, to modify/add request headers being dispatched to the upstream.
	//
	// Types that are assignable to HttpResponse:
	//
	//	*CheckResponse_DeniedResponse
	//	*CheckResponse_OkResponse
	OkHttpResponse *authv3.OkHttpResponse `cel:"http_response"`
	// Optional response metadata that will be emitted as dynamic metadata to be consumed by the next
	// filter. This metadata lives in a namespace specified by the canonical name of extension filter
	// that requires it:
	//
	// - :ref:`envoy.filters.http.ext_authz <config_http_filters_ext_authz_dynamic_metadata>` for HTTP filter.
	// - :ref:`envoy.filters.network.ext_authz <config_network_filters_ext_authz_dynamic_metadata>` for network filter.
	DynamicMetadata *structpb.Struct `cel:"dynamic_metadata"`
}

func (r OkResponse) ToCheckResponse() *authv3.CheckResponse {
	return &authv3.CheckResponse{
		Status: r.Status,
		HttpResponse: &authv3.CheckResponse_OkResponse{
			OkResponse: r.OkHttpResponse,
		},
		DynamicMetadata: r.DynamicMetadata,
	}
}

type DeniedResponse struct {
	// Status “OK“ allows the request. Any other status indicates the request should be denied, and
	// for HTTP filter, if not overridden by :ref:`denied HTTP response status <envoy_v3_api_field_service.auth.v3.DeniedHttpResponse.status>`
	// Envoy sends “403 Forbidden“ HTTP status code by default.
	Status *status.Status `cel:"status"`
	// An message that contains HTTP response attributes. This message is
	// used when the authorization service needs to send custom responses to the
	// downstream client or, to modify/add request headers being dispatched to the upstream.
	//
	// Types that are assignable to HttpResponse:
	//
	//	*CheckResponse_DeniedResponse
	//	*CheckResponse_OkResponse
	DeniedHttpResponse *authv3.DeniedHttpResponse `cel:"http_response"`
	// Optional response metadata that will be emitted as dynamic metadata to be consumed by the next
	// filter. This metadata lives in a namespace specified by the canonical name of extension filter
	// that requires it:
	//
	// - :ref:`envoy.filters.http.ext_authz <config_http_filters_ext_authz_dynamic_metadata>` for HTTP filter.
	// - :ref:`envoy.filters.network.ext_authz <config_network_filters_ext_authz_dynamic_metadata>` for network filter.
	DynamicMetadata *structpb.Struct `cel:"dynamic_metadata"`
}

func (r DeniedResponse) ToCheckResponse() *authv3.CheckResponse {
	return &authv3.CheckResponse{
		Status: r.Status,
		HttpResponse: &authv3.CheckResponse_DeniedResponse{
			DeniedResponse: r.DeniedHttpResponse,
		},
		DynamicMetadata: r.DynamicMetadata,
	}
}
