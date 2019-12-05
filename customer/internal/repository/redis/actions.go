package rediscache

import (
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
	grpcProto "github.com/shanehowearth/bcg/customer/integration/grpc/proto/v1"
)

// tmpStruct because I do not want to mess with the generated grpc tags
type tmpStruct struct {
	ID      string `redis:"id"`
	Name    string `redis:"name"`
	Address string `redis:"address"`
	Email   string `redis:"email"`
	Phone   string `redis:"phone"`
}

// Create -
func (r *Redis) Create(customer *grpcProto.CustomerDetails) (string, error) {
	// get conn and put back when exit from method
	var conn redis.Conn
	if r.Pool == nil {
		r.initPool()
		conn = r.Pool.Get()
		r.ping(conn)
	}
	if conn == nil {
		conn = r.Pool.Get()
	}
	defer conn.Close()

	if customer == nil {
		return "", fmt.Errorf("no Customer Details supplied")
	}
	err := conn.Send("INCR", "customers")
	if err != nil {
		log.Printf("actions incr %v", err)
		return "", fmt.Errorf("unable to incr with error %v", err)
	}
	tmp, err := conn.Do("GET", "customers")
	if err != nil {
		log.Printf("actions get %v", err)
		return "", fmt.Errorf("unable to get customers with error %v", err)
	}
	err = conn.Send("MULTI")
	if err != nil {
		log.Printf("actions multi %v", err)
		return "", fmt.Errorf("unable to insert %v with error %v", customer, err)
	}
	customer.Id = string(tmp.([]uint8))
	err = conn.Send("HSET", customer.Id, "name", customer.Name, "address", customer.Address, "email", customer.Email, "phone", customer.Phone)
	if err != nil {
		log.Printf("actions hset %v", err)
		return "", fmt.Errorf("unable to hset %v with error %v", customer, err)
	}
	_, err = conn.Do("EXEC")
	if err != nil {
		log.Printf("actions exec %v", err)
		return "", fmt.Errorf("unable to insert %v with error %v", customer, err)
	}
	return customer.Id, nil
}

// GetByID -
func (r *Redis) GetByID(id string) (*grpcProto.CustomerDetails, bool) {
	// get conn and put back when exit from method
	var conn redis.Conn
	if r.Pool == nil {
		r.initPool()
		conn = r.Pool.Get()
		r.ping(conn)
	}
	if conn == nil {
		conn = r.Pool.Get()
	}
	defer conn.Close()

	dataset, err := redis.Values(conn.Do("HGETALL", id))
	if err != nil {
		log.Printf("ERROR: failed get key %s, error %s", id, err.Error())
		return &grpcProto.CustomerDetails{}, false
	}

	// Put dataset into an Customer
	f := tmpStruct{}
	customer := &grpcProto.CustomerDetails{}

	if len(dataset) == 0 {
		return customer, false
	}
	err = redis.ScanStruct(dataset, &f)
	if err != nil {
		log.Printf("error scanning struct: %v", err)
	}
	customer.Id = id
	customer.Name = f.Name
	customer.Email = f.Email
	customer.Address = f.Address
	customer.Phone = f.Phone
	return customer, true
}
