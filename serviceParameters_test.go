package main

import (
	"os"
	"testing"
)

func SetUpParseServiceParameter() func() {
	prev := map[string]string{}

	for _, name := range []string{"PORT", "ADDRESS"} {
		prev[name] = os.Getenv(name)
	}

	return func() {
		for name, value := range prev {
			os.Setenv(name, value)
		}
	}
}

func TestParseServiceParameter(t *testing.T) {
	defer SetUpParseServiceParameter()()

	testSuites := []*struct {
		address      string
		port         string
		expectedAddr string
	}{
		{expectedAddr: "localhost:9090"},
		{address: "10.0.0.1", port: "7777", expectedAddr: "10.0.0.1:7777"},
		{address: "10.0.0.1", port: "invalid7777", expectedAddr: "10.0.0.1:9090"},
	}

	for _, test := range testSuites {
		os.Setenv("ADDRESS", test.address)
		os.Setenv("PORT", test.port)

		params := ParseServiceParameter()

		if exist := params.Address(); exist != test.expectedAddr {
			t.Error("failed to parse environment params: address=", test.address, "; port=", test.port, ". Got '", exist, "', but expected is ", test.expectedAddr)
		}
	}
}
