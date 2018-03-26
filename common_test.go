package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func init() {
	// switch off logging for testing
	log.SetOutput(bytes.NewBufferString(""))
}

type processedResult struct {
	sync.RWMutex

	result map[string]int
}

func NewProcessedResult() *processedResult {
	return &processedResult{
		result: make(map[string]int),
	}
}

func (o *processedResult) Add(email string) {
	o.Lock()
	defer o.Unlock()

	o.result[email] += 1
}

func (o *processedResult) Get(email string) int {
	o.RLock()
	defer o.RUnlock()

	return o.result[email]
}

type botHttpResponse struct {
	HttpStatus  int
	Timeout     time.Duration
	Paid        bool
	Break       bool
	InvalidJson bool
}

func makeTestServer(t *testing.T, responseSequences []*botHttpResponse, result *processedResult) *httptest.Server {
	var index int32 = -1

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		i := int(atomic.AddInt32(&index, 1))

		var record Record

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error("failed to read body of the test request", err)
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		err = json.Unmarshal(body, &record)
		if err != nil {
			t.Error("failed to unmarshal a request body as a Record", err)
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		if result != nil {
			result.Add(record.Email)
		}

		if i < len(responseSequences) {
			currentResponse := responseSequences[i]

			if currentResponse.Break {
				return
			}

			if currentResponse.Timeout > 0 {
				time.Sleep(currentResponse.Timeout)
			}

			if currentResponse.InvalidJson {
				body = []byte("{paid:")
			} else {
				body, _ = json.Marshal(&struct {
					Paid bool `json:"paid"`
				}{
					Paid: currentResponse.Paid,
				})
			}

			w.WriteHeader(currentResponse.HttpStatus)
			w.Write(body)
		} else {
			t.Error("got ", index+1, " request(s), but expect only ", len(responseSequences))
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
}
