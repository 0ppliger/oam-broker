package main

import (
	"log"
	"net"
	"context"
	"net/netip"

	"google.golang.org/grpc"
	"github.com/0ppliger/open-asset-gateway/api"
	"google.golang.org/protobuf/types/known/timestamppb"
	neo_db "github.com/owasp-amass/asset-db/repository/neo4j"
	network "github.com/owasp-amass/open-asset-model/network"
)

type Server struct {
	api.UnimplementedStoreServiceServer
}

func (s *Server) CreateIPAddress(
	ctx context.Context,
	input *api.IPAddressAsset,
) (*api.IPAddressEntity, error) {
	store, err := neo_db.New(neo_db.Neo4j, "bolt://neo4j:password@localhost:7687/neo4j")
	if err != nil {
		log.Fatal(err)
	}

	ip_address, err := netip.ParseAddr(input.Address)
	if err != nil {
		log.Fatal(err)
	}

	ip_type := input.Type
	if ip_type != "IPv4" && ip_type != "IPv6" {
		log.Fatalf("Wrong IP type (IPv4 or IPv6)")
	}
	

	entity, err := store.CreateAsset(&network.IPAddress{
		Address: ip_address,
		Type: ip_type,
	})
	if err != nil {
		log.Fatal(err)
	}
	
	return &api.IPAddressEntity{
		Id: entity.ID,
		Asset: input,
		CreatedAt: timestamppb.New(entity.CreatedAt),
		LastSeen: timestamppb.New(entity.LastSeen),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %v", err)
	}

	grpcServer := grpc.NewServer()

	api.RegisterStoreServiceServer(grpcServer, &Server{})

	log.Printf("Starting to listen on port 9000...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 9000: %v", err)
	}
}
