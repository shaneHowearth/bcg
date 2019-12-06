package main

import (
	"context"
	"log"
	"net"
	"os"

	grpcProto "github.com/shanehowearth/bcg/notify/integration/grpc/proto/v1"
	gmail "github.com/shanehowearth/bcg/notify/internal/emailprovider/gmail"
	"github.com/shanehowearth/bcg/notify/internal/notify"
	twilio "github.com/shanehowearth/bcg/notify/internal/smsprovider/Twilio"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var ns *notify.Server

type server struct{}

func main() {

	// SMS Setup
	sms := twilio.NewClient()

	// Mail Setup
	mail := gmail.Client{}

	// Article Service
	ns = notify.NewNotifyService(sms, mail)

	// gRPC service
	portNum := os.Getenv("PORT_NUM")
	lis, err := net.Listen("tcp", "0.0.0.0:"+portNum)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	grpcProto.RegisterNotifyServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *server) CreateNotification(ctx context.Context, cd *grpcProto.CustomerDetails) (*grpcProto.Acknowledgement, error) {

	// Send SMS notification
	err := ns.SendSMS(ctx, cd)
	if err != nil {
		log.Printf("Error sending SMS notification: %v", err)
	}
	log.Printf("SMS Sending completed")

	// Send Email
	err = ns.SendMail(ctx, cd)
	if err != nil {
		log.Printf("Error sending Mail notification: %v", err)
	}
	log.Printf("Mail Sending completed")
	return &grpcProto.Acknowledgement{}, err
}
