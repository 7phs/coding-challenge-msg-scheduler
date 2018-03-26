package main

import (
	"net/http"
	"testing"
	"time"
)

func TestSendMessageRequest_Error(t *testing.T) {
	rec := NewRecord("test@test.com", "test", nil)

	// check request witout a server
	status, err := SendMessageRequest("http://unknown.local.host:13319/", rec, 0)
	expectedStatus := messages_STATUS_ERROR

	if status != expectedStatus {
		t.Error("failed to catch an error status without a server. Got status ", status, ", but expected is ", expectedStatus)
	}
	if err == nil {
		t.Error("failed to catch an error without a server")
	}

	httpResponses := []*botHttpResponse{
		{HttpStatus: http.StatusCreated, Timeout: 100 * time.Millisecond},
		{HttpStatus: http.StatusNotFound},
		{HttpStatus: http.StatusCreated, Break: true},
		{HttpStatus: http.StatusCreated, InvalidJson: true},
	}
	testServer := makeTestServer(t, httpResponses, nil)
	defer testServer.Close()

	testSuites := []*struct {
		name    string
		timeout time.Duration
	}{
		{name: "without a server", timeout: 1 * time.Millisecond},
		{name: "invalid status"},
		{name: "break connection", timeout: 100 * time.Millisecond},
		{name: "invalid json"},
	}

	for _, test := range testSuites {
		status, err = SendMessageRequest(testServer.URL, rec, test.timeout)
		expectedStatus = messages_STATUS_ERROR

		if status != expectedStatus {
			t.Error("failed to catch an error status: ", test.name, ". Got status ", status, ", but expected is ", expectedStatus)
		}
		if err == nil {
			t.Error("failed to catch an error: ", test.name)
		}
	}
}

func TestSendMessageRequest(t *testing.T) {
	rec := NewRecord("test@test.com", "test", nil)

	httpResponses := []*botHttpResponse{
		{HttpStatus: http.StatusCreated},
		{HttpStatus: http.StatusCreated, Paid: true},
	}

	testServer := makeTestServer(t, httpResponses, nil)
	defer testServer.Close()

	testSuites := []Status{
		messages_STATUS_CONTINUE,
		messages_STATUS_COMPLETE,
	}

	for _, expectedStatus := range testSuites {
		status, err := SendMessageRequest(testServer.URL, rec, 0)

		if status != expectedStatus {
			t.Error("failed to execute a request. Got status ", expectedStatus, ", but expected is ", expectedStatus)
		}

		if err != nil {
			t.Error("execute a request with an error: ", err)
		}

	}
}
