package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	// Import your db package and protobuf definitions
	"go_tutoriols/VectorEngine/db"
	"go_tutoriols/VectorEngine/db/pb"
)

type server struct {
	pb.UnimplementedVectorServiceServer
	engine *db.VectorEngine // Use the db. prefix here
}

func (s *server) Insert(ctx context.Context, req *pb.InsertRequest) (*pb.InsertResponse, error) {
	s.engine.Insert(req.Vector)
	return &pb.InsertResponse{Success: true}, nil
}

func (s *server) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	nearest, dist := s.engine.Search(req.Query)
	return &pb.SearchResponse{
		NearestVector: nearest,
		Distance:      dist,
	}, nil
}

func main() {
	// Call NewVectorEngine from your db package
	engine := db.NewVectorEngine(10000)

	grpcServer := grpc.NewServer()
	pb.RegisterVectorServiceServer(grpcServer, &server{engine: engine})

	port := ":50051"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	fmt.Printf("Vector DB gRPC Server listening on %s...\n", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
