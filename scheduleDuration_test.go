package main

import (
	"reflect"
	"testing"
	"time"
)

func TestParseScheduleDuration(t *testing.T) {
	testSuites := []*struct {
		str           string
		expectedErr   bool
		expectedSched []time.Duration
	}{
		{str: "3f-7s-18m", expectedErr: true},
		{str: "3s-a7s-18m", expectedErr: true},
		{str: "--", expectedSched: []time.Duration{}},
		{str: "", expectedSched: []time.Duration{}},
		{str: "0s", expectedSched: []time.Duration{0 * time.Second}},
		{str: "8s-14m-20h", expectedSched: []time.Duration{
			8 * time.Second, 14 * time.Minute, 20 * time.Hour,
		}},
	}

	for _, test := range testSuites {
		schedule, err := ParseScheduleDuration(test.str)
		if test.expectedErr {
			if err == nil {
				t.Error("failed to catch error while parse schedule duration '", test.str, "'")
			}
		} else if err != nil {
			t.Error("failed to parse schedule duration with error ", err)
		} else if !reflect.DeepEqual(schedule, test.expectedSched) {
			t.Error("failed to parse schedule duration. Got ", schedule, ", but expected is ", test.expectedSched)

		}
	}
}
