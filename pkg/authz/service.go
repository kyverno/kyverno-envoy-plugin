package authz

import (
	"context"
	"encoding/json"
	"fmt"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/encoding/protojson"
)

type service struct{}

func (s *service) Check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	// parse request
	resource, err := convert(req.GetAttributes())
	if err != nil {
		return nil, err
	}
	fmt.Println(resource)
	// evaluate policies
	if r := allow(); r != nil {
		return r, nil
	}
	return deny("foo"), nil
}

// convert takes an AttributeContext and marshals it into an unstructured map
func convert(attrs *authv3.AttributeContext) (map[string]any, error) {
	// create a new Marshaler
	marshaler := protojson.MarshalOptions{
		Multiline:    true,
		Indent:       "  ",
		AllowPartial: true,
	}
	// marshal the AttributeContext to json
	jsonData, err := marshaler.Marshal(attrs)
	if err != nil {
		return nil, fmt.Errorf("maarshaling attributes to json failed: %w", err)
	}
	// unmarshal the json into an any
	var resource map[string]any
	if err := json.Unmarshal(jsonData, &resource); err != nil {
		return nil, err
	}
	return resource, nil
}

func allow() *authv3.CheckResponse {
	return &authv3.CheckResponse{
		Status:          &status.Status{Code: int32(codes.OK)},
		HttpResponse:    &authv3.CheckResponse_OkResponse{},
		DynamicMetadata: nil,
	}
}

func deny(denialReason string) *authv3.CheckResponse {
	return &authv3.CheckResponse{
		Status: &status.Status{
			Code: int32(codes.PermissionDenied),
		},
		HttpResponse: &authv3.CheckResponse_DeniedResponse{
			DeniedResponse: &authv3.DeniedHttpResponse{
				Status: &typev3.HttpStatus{Code: typev3.StatusCode_Forbidden},
				Body:   fmt.Sprintf("Request denied by Kyverno JSON engine. Reason: %s", denialReason),
			},
		},
		DynamicMetadata: nil,
	}
}
