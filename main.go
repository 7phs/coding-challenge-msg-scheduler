package main

import (
	"log"
	"net/http"
	"time"
)

func checkService(address string) error {
	_, err := (&http.Client{
		Timeout: 100 * time.Millisecond,
	}).Head("http://" + address)

	return err
}

func main() {
	log.Print("test-util starting...")

	serviceParameters := ParseServiceParameter()
	commandArgs, err := ParseCommandArgs()
	if err != nil {
		log.Fatal("failed to parse arguments: ", err)
	}

	dataSource, err := NewDataSource(commandArgs.DataFile())
	if err != nil {
		log.Fatal("failed to create a data sources: ", err)
	}

	if err = checkService(serviceParameters.Address()); err != nil {
		log.Fatal("failed to send request to a data service ", serviceParameters.Address(), ". You should run it first")
	}

	log.Print("service address: ", serviceParameters.Address())
	log.Print("data file: ", commandArgs.DataFile())

	NewMessageProcessor(serviceParameters, dataSource).Start()

	log.Print("test-util finish")
}
