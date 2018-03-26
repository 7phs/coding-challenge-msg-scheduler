package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	DEFAULT_ADDRESS = "localhost"
	DEFAULT_PORT    = 9090
)

type ServiceParameters struct {
	address string
	port    int
}

func (o *ServiceParameters) Address() string {
	return fmt.Sprintf("%s:%d", o.address, o.port)
}

func ParseServiceParameter() *ServiceParameters {
	var (
		address = os.Getenv("ADDRESS")
		portStr = os.Getenv("PORT")
		port    = DEFAULT_PORT
	)

	if address == "" {
		address = DEFAULT_ADDRESS
	}

	if p, err := strconv.Atoi(portStr); err == nil {
		port = p
	} else if portStr != "" && err != nil {
		log.Print("using default value for parameter PORT. An environment value has an error: " + err.Error())
	}

	return &ServiceParameters{
		address: address,
		port:    port,
	}
}
