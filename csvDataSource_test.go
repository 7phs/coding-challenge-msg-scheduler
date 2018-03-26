package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"reflect"
	"testing"
	"time"
)

func TestNewCsvDataSources(t *testing.T) {
	testSuites := []*struct {
		data        string
		expectedErr bool
	}{
		{data: "", expectedErr: true},
		{data: "unknown, unknown, unknown", expectedErr: true},
		{data: "email,text,schedule,text,email", expectedErr: true},
		{data: "text,schedule", expectedErr: true},
		{data: "email,text,text", expectedErr: true},
		{data: "email,text,schedule"},
		{data: "email,   text   ,   schedule   "},
	}

	for _, test := range testSuites {
		dataSource, err := NewCsvDataSources(ioutil.NopCloser(bytes.NewReader([]byte(test.data))))
		if test.expectedErr {
			if err == nil {
				t.Error("failed to catch error while parse header")
			}
		} else if err != nil {
			t.Error("failed to parse header '", test.data, "' with error ", err)
		}

		if dataSource != nil {
			dataSource.Close()
		}
	}
}

func TestNewCsvDataSources_Next(t *testing.T) {
	data := `email,text,schedule
vdaybell0@seattletimes.com,"Hi Vincenty, your invoice about $1.99 is due.",8s-14s-20s
,"Another message",0s
charrimanr@ucla.edu,,0s
charrimanr@ucla.edu,"Another message",3f-7s-18m
bskentelberyl@mozilla.org,"Dear Mr. Skentelbery, you still have an outstanding amount of $152.87 for your loan.",3s-7s-18s
`

	expected := []*Record{
		NewRecord("vdaybell0@seattletimes.com", "Hi Vincenty, your invoice about $1.99 is due.",
			[]time.Duration{8 * time.Second, 14 * time.Second, 20 * time.Second}),
		NewRecord("bskentelberyl@mozilla.org", "Dear Mr. Skentelbery, you still have an outstanding amount of $152.87 for your loan.",
			[]time.Duration{3 * time.Second, 7 * time.Second, 18 * time.Second}),
	}

	dataSource, err := NewCsvDataSources(ioutil.NopCloser(bytes.NewReader([]byte(data))))
	if err != nil {
		t.Error("failed to create a data source with error ", err)
		return
	}

	var (
		expectedIndex = 0
		existRec      *Record
	)

	for err != io.EOF {
		existRec, err = dataSource.Next()
		if err != nil || !existRec.IsValid() {
			continue
		}

		if expectedIndex < len(expected) {
			if expectedRec := expected[expectedIndex]; !reflect.DeepEqual(existRec, expectedRec) {
				t.Error("failed parse a record. Got ", existRec, ", but expected ", expectedRec)
			}

			expectedIndex++
		} else {
			t.Error("got more than ", len(expected), " records")
		}
	}

	if expectedIndex < len(expected) {
		t.Error("failed to check all expected record. Got ", expectedIndex, " record(s), but expected count is ", len(expected))
	}
}
