package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"testing"
)

func TestNewMessageProcessor(t *testing.T) {
	data := `email,text,schedule
vdaybell0@seattletimes.com,"Hi Vincenty, your invoice about $1.99 is due.",0s-0s-0s
charrimanr@ucla.edu,"Another message",0s
bskentelberyl@mozilla.org,"Dear Mr. Skentelbery, you still have an outstanding amount of $152.87 for your loan.",0s
`

	result := NewProcessedResult()

	httpResponses := []*botHttpResponse{
		{HttpStatus: http.StatusCreated},
		{HttpStatus: http.StatusCreated},
		{HttpStatus: http.StatusCreated},
		{HttpStatus: http.StatusCreated},
		{HttpStatus: http.StatusCreated},
	}
	testServer := makeTestServer(t, httpResponses, result)
	defer testServer.Close()

	dataSource, err := NewCsvDataSources(ioutil.NopCloser(bytes.NewReader([]byte(data))))
	if err != nil {
		t.Error("failed to create a data source with error ", err)
		return
	}

	testUrl, _ := url.Parse(testServer.URL)
	port, _ := strconv.Atoi(testUrl.Port())

	NewMessageProcessor(&ServiceParameters{
		address: testUrl.Hostname(),
		port:    port,
	}, dataSource).Start()

	expectedResult := map[string]int{
		"vdaybell0@seattletimes.com": 3,
		"charrimanr@ucla.edu":        1,
		"bskentelberyl@mozilla.org":  1,
	}

	for email, expectedCount := range expectedResult {
		if existCount := result.Get(email); existCount != expectedCount {
			t.Error("request counts for '", email, "' was getting ", existCount, ", but expected is ", expectedCount)
		}
	}
}
