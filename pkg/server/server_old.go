package server

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"net"
// 	"net/http"
// 	"net/url"
// 	"os"
// 	"strings"
// 	"time"

// 	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
// 	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
// 	jsonengine "github.com/kyverno/kyverno-envoy-plugin/pkg/json-engine"
// 	"github.com/kyverno/kyverno-envoy-plugin/pkg/request"
// 	"github.com/kyverno/kyverno-envoy-plugin/pkg/server/handlers"
// 	"github.com/kyverno/kyverno-envoy-plugin/pkg/signals"
// 	"github.com/kyverno/kyverno-json/pkg/policy"
// 	"go.uber.org/multierr"
// 	"google.golang.org/genproto/googleapis/rpc/status"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/codes"
// 	"k8s.io/apimachinery/pkg/util/wait"
// )

// func (s *extAuthzServerV3) Check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
// 	// Parse request
// 	attrs := req.GetAttributes()
// 	// Load policies from files
// 	policies, err := policy.Load(s.policies...)
// 	if err != nil {
// 		log.Printf("Failed to load policies: %v", err)
// 		return nil, err
// 	}

// 	resource, err := request.Convert(attrs)
// 	if err != nil {
// 		log.Printf("Error converting request: %v", err)
// 		return nil, err
// 	}

// 	engine := jsonengine.New()
// 	response := engine.Run(ctx, jsonengine.Request{
// 		Resource: resource,
// 		Policies: policies,
// 	})

// 	log.Printf("Request is initialized in kyvernojson engine .\n")

// 	var violations []error

// 	for _, policy := range response.Policies {
// 		for _, rule := range policy.Rules {
// 			if rule.Error != nil {
// 				// If there is an error, add it to the violations error array
// 				violations = append(violations, fmt.Errorf("%v", rule.Error))
// 				log.Printf("Request violation: %v\n", rule.Error.Error())
// 			} else if len(rule.Violations) != 0 {
// 				// If there are violations, add them to the violations error array
// 				for _, violation := range rule.Violations {
// 					violations = append(violations, fmt.Errorf("%v", violation))
// 				}
// 				log.Printf("Request violation: %v\n", rule.Violations.Error())
// 			} else {
// 				// If the rule passed, log it but continue to the next rule/policy
// 				log.Printf("Request passed the %v policy rule.\n", rule.Rule.Name)
// 			}
// 		}
// 	}

// 	if len(violations) == 0 {
// 		return s.allow(), nil
// 	} else {
// 		// combiine all violations errors into a single error
// 		denialReason := multierr.Combine(violations...).Error()
// 		return s.deny(denialReason), nil
// 	}

// }
