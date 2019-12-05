package customerservice

import (
	"context"
	"fmt"
	"log"
	"regexp"

	grpcProto "github.com/shanehowearth/bcg/customer/integration/grpc/proto/v1"
	repo "github.com/shanehowearth/bcg/customer/integration/repository/cache/v1"
)

// Server -
type Server struct {
	Cache repo.Cache
}

// NewCustomerService -
func NewCustomerService(c repo.Cache) *Server {
	if c == nil {
		log.Fatal("Cache supplied for NewCustomerService is nil")
	}
	a := Server{Cache: c}
	return &a
}

// CreateCustomer -
func (a *Server) CreateCustomer(ctx context.Context, det *grpcProto.CustomerDetails) (string, error) {
	log.Printf("internal/customerservice CreateCustomer %v", det)
	// Service level input validation
	if det.GetName() == "" || det.GetEmail() == "" {
		log.Printf("Missing fields name: %s email: %s", det.GetName(), det.GetEmail())
		return "", fmt.Errorf("error Name and Email are mandatory fields, please check that they have been included in your customer details")
	}
	// *Very* basic email validation
	emailRegExp := "^([a-zA-Z0-9_\\-\\.]+)@([a-zA-Z0-9_\\-\\.]+)\\.([a-zA-Z]{2,5})$"
	match, _ := regexp.MatchString(emailRegExp, det.GetEmail())
	if !match {
		log.Printf("Email %q did not match regexp", det.GetEmail())
		return "", fmt.Errorf("error Email supplied is invalid %q", det.GetEmail())
	}

	return a.Cache.Create(det)
}

// GetCustomer -
func (a *Server) GetCustomer(ctx context.Context, req *grpcProto.CustomerRequest) (*grpcProto.CustomerDetails, error) {
	log.Printf("internal/customerservice GetCustomer %v", req)
	id := req.GetId()
	customer, _ := a.Cache.GetByID(id)

	log.Printf("internal/customerservice customer: %#+v", customer)
	return customer, nil
}
