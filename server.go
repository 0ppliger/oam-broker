package main

import (
	"log"
	"fmt"
	"net"
	"errors"
	"context"

	"google.golang.org/grpc"
	"github.com/0ppliger/open-asset-gateway/api"
)

type Server struct {
	api.UnimplementedGreetingServer
}

func (s *Server) SayHello(
	ctx context.Context,
	input *api.SayHelloInput,
) (*api.SayHelloOutput, error) {

	name := input.Value
	if name == "" {
		return nil, errors.New("empty name")
	}
	
	return &api.SayHelloOutput{
		Value: fmt.Sprintf("Hello %s", name),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %v", err)
	}

	grpcServer := grpc.NewServer()

	api.RegisterGreetingServer(grpcServer, &Server{})

	log.Printf("Starting to listen on port 9000...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 9000: %v", err)
	}
}
