package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/go-chi/chi"

	cgrpcProto "github.com/shanehowearth/bcg/customer/integration/grpc/proto/v1"
)

func bcgRoutes(router *chi.Mux) {
	// Customer related routes
	router.Route("/customer", func(r chi.Router) {
		r.Post("/", CreateCustomer)
		r.Route("/{customerID}", func(r2 chi.Router) {
			r2.Get("/", GetCustomerByID)
		})
	})
}

// GetCustomerByID -
func GetCustomerByID(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "customerID")
	// validate (id can only be int32 for now)
	if _, err := strconv.Atoi(id); err != nil {
		log.Printf("An invalid customer id was supplied, ID: %s Error: %v", id, err)
		respondWithError(w, http.StatusInternalServerError, "Supplied Customer ID is an incorrect format")
	}

	customer, err := rc.GetCustomer(&cgrpcProto.CustomerRequest{Id: id})
	errStr := randSeq(6)
	if err != nil {
		log.Printf("%s An error occurred with GetCustomerByID, Error: %v", errStr, err)
		// We don't want the user to know about the inner workings of the application
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("An internal server error has occured, please contact Customer Support and quote this unique ID %s", errStr))
	}
	respondWithJSON(w, http.StatusOK, customer)
}

// CreateCustomer -
func CreateCustomer(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var customer *cgrpcProto.CustomerDetails
	err := decoder.Decode(&customer)
	errStr := randSeq(6)
	if err != nil {
		log.Printf("%s An error occurred with CreateCustomer, Error: %v", errStr, err)
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("I am unable to use the information supplied, please try again. Alternatively you may contact Customer Support and quote this unique ID %s", errStr))
		return
	}
	// Validate
	// Name and Email are Mandatory
	if customer.Name == "" || customer.Email == "" {
		log.Printf("%s Name %s or Email %s not supplied", errStr, customer.Name, customer.Email)
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Name and Email are mandatory fields, please try again. Alternatively you may contact Customer Support and quote this unique ID %s", errStr))
		return
	}
	// *Very* basic email validation
	emailRegExp := "^([a-zA-Z0-9_\\-\\.]+)@([a-zA-Z0-9_\\-\\.]+)\\.([a-zA-Z]{2,5})$"
	match, _ := regexp.MatchString(emailRegExp, customer.Email)
	if !match {
		log.Printf("%s Bad Email %s supplied", errStr, customer.Email)
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Email supplied is badly formatted, please try again. Alternatively you may contact Customer Support and quote this unique ID %s", errStr))

		return
	}
	ack, err := rc.CreateCustomer(customer)
	if err != nil {
		log.Printf("%s An error occurred with CreateCustomer, Error: %v", errStr, err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("An internal server error has occured, please contact Customer Support and quote this unique ID %s", errStr))
		return
	}
	respondWithJSON(w, http.StatusOK, ack)

}

var letters = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Generate a pseudo random string for the customer to quote to the CSR
func randSeq(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
