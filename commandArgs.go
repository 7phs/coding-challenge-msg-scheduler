package main

import (
	"errors"
	"os"
	"strings"
)

type CommandArgs struct {
	dataFile string
}

func ParseCommandArgs() (*CommandArgs, error) {
	// check arguments count
	if len(os.Args) <= 1 || len(os.Args) > 2 {
		return nil, errors.New("not enough parameter. " +
			"Usage example: test-util data.file (" +
			"supported data formats: " + strings.Join(supportedDataFormat, ";") + ")")
	}
	// check file exists
	dataFile := os.Args[1]
	if fileInfo, err := os.Stat(dataFile); err != nil || fileInfo.IsDir() {
		return nil, errors.New("a data file '" + dataFile + "' isn't exist")
	}
	// check supported data formats
	dataFormat := parseDataFormat(dataFile)
	if !isDataFormatSupported(dataFormat) {
		return nil, errors.New("a data format '" + dataFormat + "' not supported. Using another data file")
	}

	return &CommandArgs{
		dataFile: dataFile,
	}, nil
}

func (o *CommandArgs) DataFile() string {
	return o.dataFile
}
