package customerservice_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/bouk/monkey"
	grpcProto "github.com/shanehowearth/bcg/customer/integration/grpc/proto/v1"
	repo "github.com/shanehowearth/bcg/customer/integration/repository/cache/v1"
	SUT "github.com/shanehowearth/bcg/customer/internal/customerservice"
	"github.com/stretchr/testify/assert"
)

type mockRepoCache struct{}

var cacheCustomer *grpcProto.CustomerDetails
var cacheFound bool
var populateErr error

func (m *mockRepoCache) GetByID(id string) (*grpcProto.CustomerDetails, bool) {
	return cacheCustomer, cacheFound
}

var customerID string
var cacheError error

func (m *mockRepoCache) Create(*grpcProto.CustomerDetails) (string, error) {
	return customerID, cacheError
}

// Begin tests

func TestNewCustomerService(t *testing.T) {
	mockCache := &mockRepoCache{}
	testcases := map[string]struct {
		cache       repo.Cache
		server      SUT.Server
		errMessage  string
		expectPanic bool
		cacheErr    error
	}{
		"Happy Path":          {cache: mockCache, server: SUT.Server{Cache: mockCache}},
		"Missing Cache":       {expectPanic: true, errMessage: "Cache supplied for NewCustomerService is nil"},
		"Cache returns error": {cache: mockCache, server: SUT.Server{Cache: mockCache}, cacheErr: fmt.Errorf("error returned")},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			populateErr = tc.cacheErr
			if tc.expectPanic {
				fakeLogFatal := func(msg ...interface{}) {
					assert.Equal(t, tc.errMessage, msg[0])
					panic("log.Fatal called")
				}
				patch := monkey.Patch(log.Fatal, fakeLogFatal)
				defer patch.Unpatch()
				assert.PanicsWithValue(t, "log.Fatal called", func() { SUT.NewCustomerService(tc.cache) }, "log.Fatal was not called")
			} else {

				output := SUT.NewCustomerService(tc.cache)
				assert.Equal(t, *output, tc.server)
			}
		})
	}
}

func TestCreateCustomer(t *testing.T) {
	mockCache := &mockRepoCache{}

	testcases := map[string]struct {
		ctx      context.Context
		input    *grpcProto.CustomerDetails
		response string
		err      bool
	}{
		"Happy Path": {
			ctx:      context.Background(),
			input:    &grpcProto.CustomerDetails{Name: "Test", Email: "test@domain.com"},
			response: "1",
		},
		"Missing Name": {
			ctx:   context.Background(),
			input: &grpcProto.CustomerDetails{Email: "test@domain.com"},
			err:   true,
		},
		"Missing Email": {
			ctx:   context.Background(),
			input: &grpcProto.CustomerDetails{Name: "Test"},
			err:   true,
		},
		"Malformed Email": {
			ctx:   context.Background(),
			input: &grpcProto.CustomerDetails{Name: "Test", Email: "testdomain.com"},
			err:   true,
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			customerID = tc.response
			ss := SUT.NewCustomerService(mockCache)
			output, err := ss.CreateCustomer(tc.ctx, tc.input)
			assert.Equal(t, tc.response, output, "Expected %v got %v", tc.response, output)
			if tc.err {

				assert.NotNil(t, err, "Was expecting an error")
			} else {
				assert.Nil(t, err, "Not expecting an error")
			}
		})
	}
}

func TestGetCustomer(t *testing.T) {
	mockCache := &mockRepoCache{}

	testcases := map[string]struct {
		ctx           context.Context
		input         *grpcProto.CustomerRequest
		response      *grpcProto.CustomerDetails
		found         bool
		errorReturned bool
		cacheErr      error
	}{
		"Happy Path": {
			ctx:      context.Background(),
			input:    &grpcProto.CustomerRequest{Id: "1"},
			response: &grpcProto.CustomerDetails{Id: "1"},
			found:    true,
		},
		"No Customer found": {
			ctx:      context.Background(),
			input:    &grpcProto.CustomerRequest{Id: "1"},
			response: &grpcProto.CustomerDetails{},
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {

			ss := SUT.NewCustomerService(mockCache)
			cacheCustomer = tc.response
			cacheFound = tc.found
			populateErr = tc.cacheErr

			output, err := ss.GetCustomer(tc.ctx, tc.input)
			assert.Equal(t, *tc.response, *output, "Expected %v got %v", tc.response, output)
			if tc.errorReturned {
				assert.NotNil(t, err, "Expecting error")
			} else {
				assert.Nil(t, err, "Not expecting an error, but got %v", err)
			}
		})
	}
}
