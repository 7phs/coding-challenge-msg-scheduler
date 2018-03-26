package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Status int

const (
	messages_STATUS_ERROR Status = -1 - iota
)

const (
	messages_STATUS_CONTINUE Status = iota
	messages_STATUS_COMPLETE
)

const (
	messages_ENDPOINT = "/messages"
)

type MessageResponse struct {
	Paid bool `json:"paid"`
}

func SendMessageRequest(address string, rec *Record, timeout time.Duration) (Status, error) {
	body, err := json.Marshal(rec)
	if err != nil {
		return messages_STATUS_ERROR, errors.New(fmt.Sprint("failed to marshal with error ", err))
	}

	resp, err := (&http.Client{
		Timeout: timeout,
	}).Post(address+messages_ENDPOINT, "text/json", bytes.NewReader(body))
	if err != nil {
		return messages_STATUS_ERROR, errors.New(fmt.Sprint("failed to send request with error ", err))
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return messages_STATUS_ERROR, errors.New(fmt.Sprint("failed to read response body with error ", err))
	}

	if resp.StatusCode != http.StatusCreated {
		return messages_STATUS_ERROR, errors.New(fmt.Sprint("got error status=", resp.Status, "; body=", string(respBody)))
	}

	var messageResponse MessageResponse
	if err = json.Unmarshal(respBody, &messageResponse); err != nil {
		return messages_STATUS_ERROR, errors.New(fmt.Sprint("failed to unmarshal response body with error ", err))
	}

	if messageResponse.Paid {
		return messages_STATUS_COMPLETE, nil
	}

	return messages_STATUS_CONTINUE, nil
}
