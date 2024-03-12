package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
)

func startHTTPserver() {
	http.HandleFunc("/", handler)

	http.ListenAndServe(":8080", nil)
	fmt.Println("Starting HTTP server on Port 8080")
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

func startGRPCServer() {

	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	fmt.Println("Starting GRPC server on Port 9090")
	grpcServer.Serve(lis)

}

func main() {

	startHTTPserver()
	startGRPCServer()

}
