package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	// Required for gRPC telemetry hooks if using stats handler, or use interceptors
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	// Import your db package and protobuf definitions
	"go_tutoriols/VectorEngine/db"
	"go_tutoriols/VectorEngine/db/pb"
)

// Define Prometheus metrics for gRPC requests
var (
	grpcRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total number of gRPC requests received by the vector worker.",
		},
		[]string{"method"},
	)
)

func init() {
	prometheus.MustRegister(grpcRequestsTotal)
}

// Unary interceptor to automatically increment request counts per method
func prometheusInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	grpcRequestsTotal.WithLabelValues(info.FullMethod).Inc()
	return handler(ctx, req)
}

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

// GetStats returns the current vector count stored in this node's engine
func (s *server) GetStats(ctx context.Context, req *pb.StatsRequest) (*pb.StatsResponse, error) {
	// Call Size() or Count() on your db.VectorEngine instance
	vectorCount := int64(s.engine.Size())

	return &pb.StatsResponse{
		Count: vectorCount,
	}, nil
}

func main() {
	// Start a background HTTP metrics server for Prometheus / HPA scraping on port 2112
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		fmt.Println("Prometheus metrics server listening on :2112/metrics...")
		if err := http.ListenAndServe(":2112", mux); err != nil {
			log.Printf("Metrics server error: %v", err)
		}
	}()

	// Call NewVectorEngine from your db package
	engine := db.NewVectorEngine(10000, 10000)

	// Register the Prometheus unary interceptor to track incoming gRPC traffic
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(prometheusInterceptor),
	)
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
