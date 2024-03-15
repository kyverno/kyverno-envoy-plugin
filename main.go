package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

type Servers struct {
	httpServer *http.Server
	grpcServer *grpc.Server
}

func NewServers() *Servers {
	return &Servers{}
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
	fmt.Fprint(w, "Hello World!")
}

func (s *Servers) startGRPCServer(ctx context.Context) {

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s.grpcServer = grpc.NewServer()
	fmt.Println("Starting GRPC server on Port 9000")

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

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	srv := NewServers()

	go srv.startHTTPServer(ctx)
	go srv.startGRPCServer(ctx)

	go func() {
		<-ctx.Done()
		cancel()
	}()
	select {}
}
