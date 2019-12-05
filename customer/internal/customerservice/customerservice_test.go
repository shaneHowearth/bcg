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
/* var tagInfo *grpcProto.TagInfo */
var populateErr error

func (m *mockRepoCache) GetByID(id string) (*grpcProto.CustomerDetails, bool) {
	return cacheCustomer, cacheFound
}
/* func (m *mockRepoCache) GetTagInfo(tag, date string) *grpcProto.TagInfo { return tagInfo } */
func (m *mockRepoCache) Populate(*grpcProto.CustomerDetails) error          { return populateErr }

type mockStorage struct{}

var fetchLatestRowsErr error
var fetchOneError error
var fetchOneCustomer *grpcProto.CustomerDetails

func (ms *mockStorage) FetchLatestRows(n int) (as []*grpcProto.CustomerDetails, e error) {
	return as, fetchLatestRowsErr
}

func (ms *mockStorage) FetchOne(id int) (a *grpcProto.CustomerDetails, e error) {
	return fetchOneCustomer, fetchOneError
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

func TestTagInfo(t *testing.T) {
	mockCache := &mockRepoCache{}

	testcases := map[string]struct {
		ctx      context.Context
		input    *grpcProto.CustomerRequest
		/* response *grpcProto.TagInfo */
	}{
		"Happy Path": {
			ctx:      context.Background(),
			input:    &grpcProto.CustomerRequest{},
			/* response: &grpcProto.TagInfo{}, */
		},
	}
	/* for name, tc := range testcases { */
	/* 	t.Run(name, func(t *testing.T) { */
	/* 		/1* tagInfo = tc.response *1/ */
	/* 		/1* ss := SUT.NewCustomerService(mockCache) *1/ */
	/* 		/1* output, err := ss.GetTagInfo(tc.ctx, tc.input) *1/ */
            /* /1* output := "" *1/ */
	/* 		/1* assert.Equal(t, tc.response, output, "Expected %v got %v", tc.response, output) *1/ */
	/* 		/1* assert.Nil(t, err, "Not expecting an error") *1/ */
	/* 	}) */
	/* } */
}

func TestGetCustomer(t *testing.T) {
	mockCache := &mockRepoCache{}

	testcases := map[string]struct {
		ctx             context.Context
		input           *grpcProto.CustomerRequest
		response        *grpcProto.CustomerDetails
		fetchedCustomer *grpcProto.CustomerDetails
		found           bool
		errorReturned   bool
		cacheErr        error
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
		"No Customer in Cache, but one in DB": {
			ctx:             context.Background(),
			input:           &grpcProto.CustomerRequest{Id: "1"},
			response:        &grpcProto.CustomerDetails{Id: "1"},
			fetchedCustomer: &grpcProto.CustomerDetails{Id: "1"},
		},
		"Bad Id supplied": {
			ctx:             context.Background(),
			input:           &grpcProto.CustomerRequest{Id: "Bad"},
			response:        &grpcProto.CustomerDetails{},
			fetchedCustomer: &grpcProto.CustomerDetails{},
			errorReturned:   true,
		},
		"Unable to populate cache": {
			ctx:             context.Background(),
			input:           &grpcProto.CustomerRequest{Id: "22"},
			response:        &grpcProto.CustomerDetails{},
			fetchedCustomer: &grpcProto.CustomerDetails{},
			cacheErr:        fmt.Errorf("cache error"),
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {

			ss := SUT.NewCustomerService(mockCache)
			cacheCustomer = tc.response
			cacheFound = tc.found
			fetchOneCustomer = tc.fetchedCustomer
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
