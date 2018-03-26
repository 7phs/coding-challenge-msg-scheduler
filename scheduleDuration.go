package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	durationSuffix = map[string]time.Duration{
		"s": time.Second,
		"m": time.Minute,
		"h": time.Hour,
	}
)

func ParseScheduleDuration(v string) ([]time.Duration, error) {
	schedules := make([]time.Duration, 0)

	for _, v := range strings.Split(v, "-") {
		v = strings.ToLower(strings.TrimSpace(v))
		if len(v) <= 1 {
			continue
		}

		suffixIndex := len(v) - 1
		duration, ok := durationSuffix[v[suffixIndex:]]
		if !ok {
			return nil, errors.New(fmt.Sprint("failed to parse a schedule '" + v + "' has an unknown duration type '" + v[suffixIndex:] + "'"))
		}

		value, err := strconv.Atoi(v[:suffixIndex])
		if err != nil {
			return nil, errors.New(fmt.Sprint("failed to parse a schedule '" + v + "' has an error while parse a time value: " + err.Error()))
		}

		schedules = append(schedules, time.Duration(value)*duration)
	}

	return schedules, nil
}
