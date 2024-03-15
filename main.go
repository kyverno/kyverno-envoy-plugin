package main

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
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"k8s.io/apimachinery/pkg/util/wait"
)

type Servers struct {
	httpServer *http.Server
	grpcServer *grpc.Server
	grpcV3     *extAuthzServerV3
}

type (
	extAuthzServerV3 struct{}
)

func NewServers() *Servers {
	return &Servers{
		grpcV3: &extAuthzServerV3{},
	}
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

func (s *extAuthzServerV3) Check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {

	attrs := req.GetAttributes()

	// Print each attribute individually
	for key, value := range attrs.GetRequest().GetHttp().GetHeaders() {
		fmt.Printf("Header: %s = %s\n", key, value)
	}

	// Print the entire struct with field names
	fmt.Printf("Attributes: %+v\n", attrs)

	// Implement your authorization logic here
	// For now, allow all requests
	return &authv3.CheckResponse{
		Status: &status.Status{Code: int32(codes.OK)},
	}, nil
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

func main() {
	var group wait.Group
	srv := NewServers()
	func() {
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer cancel()
		group.StartWithContext(ctx, srv.startHTTPServer)
		group.StartWithContext(ctx, srv.startGRPCServer)
		<-ctx.Done()
	}()
	group.Wait()
}
