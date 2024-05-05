package server

import (
	"context"

	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	jsonengine "github.com/kyverno/kyverno-envoy-plugin/pkg/json-engine"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/request"
	"github.com/kyverno/kyverno-json/pkg/policy"
	"go.uber.org/multierr"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"k8s.io/apimachinery/pkg/util/wait"
)

type extAuthzServerV3 struct {
	policies []string
}

type Servers struct {
	httpServer *http.Server
	grpcServer *grpc.Server
	grpcV3     *extAuthzServerV3
}

func NewServers(policies []string) *Servers {
	return &Servers{
		grpcV3: &extAuthzServerV3{
			policies: policies,
		},
	}
}

func StartServers(srv *Servers) {
	var group wait.Group
	func() {
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer cancel()
		group.StartWithContext(ctx, srv.startHTTPServer)
		group.StartWithContext(ctx, srv.startGRPCServer)
		<-ctx.Done()
	}()
	group.Wait()
}

func (s *Servers) startHTTPServer(ctx context.Context) {
	s.httpServer = &http.Server{
		Addr:    ":8000",
		Handler: http.HandlerFunc(handler),
	}
	fmt.Println("Starting HTTP server on Port 8000")
	go func() {
		<-ctx.Done()

		fmt.Println("HTTP server shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.httpServer.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
			log.Fatal("Shutdown HTTP server:", err)
		}
	}()

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Received request from %s %s\n", r.RemoteAddr, r.URL.Path)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	fmt.Println("Request payload:", string(body))
}

func (s *Servers) startGRPCServer(ctx context.Context) {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s.grpcServer = grpc.NewServer()
	fmt.Println("Starting GRPC server on Port 9000")

	authv3.RegisterAuthorizationServer(s.grpcServer, s.grpcV3)

	go func() {
		<-ctx.Done()
		if s.grpcServer != nil {
			fmt.Println("GRPC server shutting down...")
			s.grpcServer.GracefulStop()
		}
	}()

	if err := s.grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *extAuthzServerV3) Check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	// Parse request
	attrs := req.GetAttributes()
	// Load policies from files
	policies, err := policy.Load(s.policies...)
	if err != nil {
		log.Printf("Failed to load policies: %v", err)
		return nil, err
	}

	resource, err := request.Convert(attrs)
	if err != nil {
		log.Printf("Error converting request: %v", err)
		return nil, err
	}

	engine := jsonengine.New()
	response := engine.Run(ctx, jsonengine.Request{
		Resource: resource,
		Policies: policies,
	})

	log.Printf("Request is initialized in kyvernojson engine .\n")

	var violations []error

	for _, policy := range response.Policies {
		for _, rule := range policy.Rules {
			if rule.Error != nil {
				// If there is an error, add it to the violations error array
				violations = append(violations, fmt.Errorf("%v", rule.Error))
				log.Printf("Request violation: %v\n", rule.Error.Error())
			} else if len(rule.Violations) != 0 {
				// If there are violations, add them to the violations error array
				for _, violation := range rule.Violations {
					violations = append(violations, fmt.Errorf("%v", violation))
				}
				log.Printf("Request violation: %v\n", rule.Violations.Error())
			} else {
				// If the rule passed, log it but continue to the next rule/policy
				log.Printf("Request passed the %v policy rule.\n", rule.Rule.Name)
			}
		}
	}

	if len(violations) == 0 {
		return s.allow(), nil
	} else {
		// combiine all violations errors into a single error
		denialReason := multierr.Combine(violations...).Error()
		return s.deny(denialReason), nil
	}

}

func (s *extAuthzServerV3) allow() *authv3.CheckResponse {
	return &authv3.CheckResponse{
		Status: &status.Status{Code: int32(codes.OK)},
	}
}

func (s *extAuthzServerV3) deny(denialReason string) *authv3.CheckResponse {
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
	}
}
