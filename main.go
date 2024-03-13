package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
)

type Servers struct {
	httpServer *http.Server
	grpcServer *grpc.Server
}

func NewServers() *Servers {
	return &Servers{}
}

func (s *Servers) startHTTPServer() {

	s.httpServer = &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(handler),
	}
	fmt.Println("Starting HTTP server on Port 8080")
	if err := s.httpServer.ListenAndServe(); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

func (s *Servers) startGRPCServer() {

	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s.grpcServer = grpc.NewServer()
	fmt.Println("Starting GRPC server on Port 9090")
	if err := s.grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {

	srv := NewServers()
	go srv.startHTTPServer()
	go srv.startGRPCServer()
	select {}

}
