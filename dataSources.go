package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	for ext := range dataSourcesFabric {
		supportedDataFormat = append(supportedDataFormat, ext)
	}
}

type DataSource interface {
	Next() (*Record, error)
	Close()
}

var (
	dataSourcesFabric = map[string]func(io.ReadCloser) (DataSource, error){
		".csv": NewCsvDataSources,
	}
	supportedDataFormat []string
)

func NewDataSource(fileName string) (DataSource, error) {
	dataFormat := parseDataFormat(fileName)

	if !isDataFormatSupported(dataFormat) {
		return nil, errors.New("data format '" + dataFormat + "' isn't support")
	}

	dataFile, err := os.Open(fileName)
	if err != nil {
		return nil, errors.New("failed to open data file '" + fileName + "' with error: " + err.Error())
	}

	return dataSourcesFabric[dataFormat](dataFile)
}

func parseDataFormat(fileName string) string {
	return strings.ToLower(filepath.Ext(fileName))
}

func isDataFormatSupported(dataFormat string) bool {
	_, ok := dataSourcesFabric[dataFormat]

	return ok
}
