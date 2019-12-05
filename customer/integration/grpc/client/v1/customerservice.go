package readarticleclient

import (
	"context"
	"log"
	"time"

	grpcProto "github.com/shanehowearth/bcg/customer/integration/grpc/proto/v1"
	"google.golang.org/grpc"
)

// CustomerClient -
type CustomerClient struct {
	Address string
}

func (s *CustomerClient) newConnection() (grpcProto.CustomerServiceClient, *grpc.ClientConn) {

	// Set up a connection to the server.
	conn, err := grpc.Dial(s.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return grpcProto.NewCustomerServiceClient(conn), conn
}

// GetCustomer -
func (s *CustomerClient) GetCustomer(cr *grpcProto.CustomerRequest) (*grpcProto.CustomerDetails, error) {
	c, conn := s.newConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.GetCustomer(ctx, cr)
}

// CreateCustomer
func (s *CustomerClient) CreateCustomer(cd *grpcProto.CustomerDetails) (*grpcProto.Acknowledgement, error) {
	c, conn := s.newConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.CreateCustomer(ctx, cd)
}
