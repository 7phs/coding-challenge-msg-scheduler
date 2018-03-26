package main

import "testing"

func TestNewDataSource(t *testing.T) {
	testSuites := []*struct {
		dataFile    string
		expectedErr bool
	}{
		{dataFile: "unknown.csv", expectedErr: true},
		{dataFile: "test-data/unknown.format", expectedErr: true},
		{dataFile: "test-data/customers.csv"},
	}

	for _, test := range testSuites {
		dataSource, err := NewDataSource(test.dataFile)
		if test.expectedErr {
			if err==nil {
				t.Error("failed to catch an error for a data file '", test.dataFile, "'")
			}
		} else if err!=nil {
			t.Error("failed to create a data source for a data file '", test.dataFile, "' with error ", err)
		} else if dataSource==nil {
			t.Error("failed to create a data source for a data file '", test.dataFile, "' without any error")
		}
	}
}
