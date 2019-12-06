package messenger

import grpcProto "github.com/shanehowearth/bcg/notify/integration/grpc/proto/v1"

type Messenger interface {
	Send(recipient *grpcProto.CustomerDetails, message string) error
}
