package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	csv_FIELD_EMAIL = iota
	csv_FIELD_TEXT
	csv_FIELD_SCHEDULE

	csv_EXPECTED_FIELD_COUNT = 3
)

var (
	csvFieldMapper = map[string]int{
		"email":    csv_FIELD_EMAIL,
		"text":     csv_FIELD_TEXT,
		"schedule": csv_FIELD_SCHEDULE,
	}
)

type fieldsChecker map[int]bool

func newFieldsCheck() fieldsChecker {
	return map[int]bool{
		csv_FIELD_EMAIL:    false,
		csv_FIELD_TEXT:     false,
		csv_FIELD_SCHEDULE: false,
	}
}

func (o *fieldsChecker) Add(index int) {
	ref := (*map[int]bool)(o)

	(*ref)[index] = true
}

func (o fieldsChecker) IsComplete() bool {
	fieldsCount := 0

	for _, v := range o {
		if v {
			fieldsCount++
		}
	}

	return fieldsCount == csv_EXPECTED_FIELD_COUNT
}

type CsvDataSources struct {
	closer io.Closer
	reader *csv.Reader

	currentLine int

	fieldsMapper map[int]int
}

func NewCsvDataSources(data io.ReadCloser) (DataSource, error) {
	return (&CsvDataSources{
		closer:       data,
		reader:       csv.NewReader(data),
		fieldsMapper: make(map[int]int),
	}).parseHeader()
}

func (o *CsvDataSources) read() ([]string, error) {
	values, err := o.reader.Read()
	if err != nil {
		o.currentLine++
	}

	return values, err
}

func (o *CsvDataSources) parseHeader() (DataSource, error) {
	headers, err := o.read()
	if err != nil {
		return nil, errors.New("failed to parse a header of data file with error " + err.Error())
	}

	// check for all requirement field exists
	fieldsCheck := newFieldsCheck()
	// check duplicated, etc.
	fieldsCount := 0

	for i, name := range headers {
		name = strings.ToLower(strings.TrimSpace(name))

		index, ok := csvFieldMapper[name]
		if !ok {
			continue
		}

		fieldsCount++

		o.fieldsMapper[i] = index
		fieldsCheck.Add(index)
	}

	if fieldsCount != csv_EXPECTED_FIELD_COUNT {
		return nil, errors.New(fmt.Sprintf("failed to parse a header of data file has %d field(s), but expected %d", fieldsCount, csv_EXPECTED_FIELD_COUNT))
	}

	if !fieldsCheck.IsComplete() {
		return nil, errors.New(fmt.Sprintf("failed to parse a header of data file, not all required fields exists"))
	}

	return o, nil
}

func (o *CsvDataSources) Next() (*Record, error) {
	values, err := o.read()
	if err != nil {
		if err != io.EOF {
			err = errors.New(fmt.Sprintf("failed to read a line #%d of data file with error %s", o.currentLine, err.Error()))
		}

		return nil, err
	}

	schedules, err := ParseScheduleDuration(o.getValue(values, csv_FIELD_SCHEDULE))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to parse a schedule field of line #%d with error %s", o.currentLine, err.Error()))
	}

	email := o.getValue(values, csv_FIELD_EMAIL)
	text := o.getValue(values, csv_FIELD_TEXT)

	return NewRecord(email, text, schedules), nil
}

func (o *CsvDataSources) getValue(values []string, index int) string {
	return strings.TrimSpace(values[o.fieldsMapper[index]])
}

func (o *CsvDataSources) Close() {
	o.closer.Close()
}
