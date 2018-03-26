package main

import (
	"os"
	"strings"
	"testing"
)

func SetUpParseCommandArgs() func() {
	prev := os.Args

	return func() {
		os.Args = prev
	}
}

func TestParseCommandArgs(t *testing.T) {
	defer SetUpParseCommandArgs()()

	testSuites := []*struct {
		args             []string
		expectedErr      bool
		expectedDataFile string
	}{
		{args: []string{"test-util"},
			expectedErr: true, expectedDataFile: ""},
		{args: []string{"test-util", "data1.csv", "data2.csv"},
			expectedErr: true, expectedDataFile: ""},
		{args: []string{"test-util", "unknown.csv"},
			expectedErr: true, expectedDataFile: ""},
		{args: []string{"test-util", "test-data/customers.csv"},
			expectedDataFile: "test-data/customers.csv"},
		{args: []string{"test-util", "test-data/unknown.format"},
			expectedErr: true, expectedDataFile: ""},
	}

	for _, test := range testSuites {
		os.Args = test.args

		commandArgs, err := ParseCommandArgs()
		if test.expectedErr {
			if err == nil {
				t.Error("failed to catch an error for args: ", strings.Join(test.args, " "))
			}
		} else {
			if err != nil {
				t.Error("failed to parse command args with errpr: ", err)
			} else if exist := commandArgs.DataFile(); exist != test.expectedDataFile {
				t.Error("failed to parse command args. Got data file name='", exist, "', but expected is '", test.expectedDataFile, "'")
			}
		}
	}
}
