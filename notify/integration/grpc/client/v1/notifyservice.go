package readarticleclient

import (
	"context"
	"log"
	"time"

	grpcProto "github.com/shanehowearth/bcg/notify/integration/grpc/proto/v1"
	"google.golang.org/grpc"
)

// NotifyClient -
type NotifyClient struct {
	Address string
}

func (s *NotifyClient) newConnection() (grpcProto.NotifyServiceClient, *grpc.ClientConn) {

	// Set up a connection to the notify server.
	conn, err := grpc.Dial(s.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return grpcProto.NewNotifyServiceClient(conn), conn
}

// CreateNotification -
func (s *NotifyClient) CreateNotification(cd *grpcProto.CustomerDetails) (*grpcProto.Acknowledgement, error) {
	c, conn := s.newConnection()
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.CreateNotification(ctx, cd)
}
