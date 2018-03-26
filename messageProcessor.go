package main

import (
	"io"
	"log"
	"sync"
	"time"
)

type MessageProcessor struct {
	address    string
	dataSource DataSource
}

func NewMessageProcessor(parameters *ServiceParameters, dataSource DataSource) *MessageProcessor {
	return &MessageProcessor{
		address:    "http://" + parameters.Address(),
		dataSource: dataSource,
	}
}

func (o *MessageProcessor) Start() {
	var (
		rec   *Record
		err   error
		wait  sync.WaitGroup
		start = time.Now()
	)

	log.Print("start message processing")

	for err != io.EOF {
		rec, err = o.dataSource.Next()
		if err != nil {
			if err != io.EOF {
				log.Print("error parsing a line: ", err)
			}
			continue
		}

		if !rec.IsValid() {
			log.Print("parsed a rec is invalid: ", rec)
			continue
		}

		wait.Add(1)
		go func(rec *Record) {
			for {
				nextDuration, ok := rec.NextDuration(start)
				if !ok {
					log.Print("send message ", rec, ": failed", rec.TryingString())
					break
				} else if nextDuration > 0 {
					time.Sleep(nextDuration)
				}

				log.Print("send message ", rec, rec.TryingString())

				status, err := SendMessageRequest(o.address, rec, 0)
				if err != nil {
					log.Print("failed to send message ", rec, " with error ", err)
				} else if status == messages_STATUS_COMPLETE {
					log.Print("send message ", rec, ": success")
					break
				} else if status == messages_STATUS_CONTINUE {
					log.Print("send message ", rec, ": not complete, try again")
				}
			}

			wait.Done()
		}(rec)
	}

	wait.Wait()

	log.Print("finish message processing")
}
