package main

import (
	"context"
	"log"
	"net"
	"os"

	customer "github.com/shanehowearth/bcg/customer/internal/customerservice"
	repo "github.com/shanehowearth/bcg/customer/internal/repository/redis"

	grpcProto "github.com/shanehowearth/bcg/customer/integration/grpc/proto/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

var ss = customer.NewCustomerService(new(repo.Redis))

func main() {

	portNum := os.Getenv("PORT_NUM")
	lis, err := net.Listen("tcp", "0.0.0.0:"+portNum)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	grpcProto.RegisterCustomerServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// CreateCustomer -
func (s *server) CreateCustomer(ctx context.Context, req *grpcProto.CustomerDetails) (*grpcProto.Acknowledgement, error) {
	log.Printf("cmd/main Create Customer %v", req)
	id, err := ss.CreateCustomer(ctx, req)
	return &grpcProto.Acknowledgement{Id: id}, err
}

// GetCustomer -
func (s *server) GetCustomer(ctx context.Context, req *grpcProto.CustomerRequest) (*grpcProto.CustomerDetails, error) {
	log.Printf("cmd/main Get Customer %v", req)
	st, err := ss.GetCustomer(ctx, req)
	if err != nil {
		return nil, err
	}
	log.Printf("cmd/main customer %#+v", st)
	return st, nil
}
