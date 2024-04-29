package request

import (
	"encoding/json"
	"fmt"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/protobuf/encoding/protojson"
)

// Convert takes an AttributeContext and returns it as an any.
// It marshals the AttributeContext to JSON and then unmarshals it into an any.
func Convert(attrs *authv3.AttributeContext) (any, error) {
	// Create a new Marshaler
	marshaler := protojson.MarshalOptions{
		Multiline:    true,
		Indent:       "  ",
		AllowPartial: true,
	}

	// Marshal the AttributeContext to JSON
	jsonData, err := marshaler.Marshal(attrs)
	if err != nil {
		return nil, fmt.Errorf("maarshaling attributes to json failed: %w", err)
	}

	// Unmarshal the JSON into an any
	var resource any
	if err := json.Unmarshal(jsonData, &resource); err != nil {
		return nil, err
	}
	return resource, nil
}
