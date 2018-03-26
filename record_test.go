package main

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestNewRecord(t *testing.T) {
	if rec := NewRecord("", "", nil); rec.IsValid() {
		t.Error("failed to catch an invalid record")
	}

	if rec := NewRecord("test@test.test", "message", nil); !rec.IsValid() {
		t.Error("failed to validate a valid record")
	}
}

func TestRecord_String(t *testing.T) {
	var (
		email = "test@test.test"
		text  = "test"
		rec   = NewRecord(email, text, nil)
	)

	str := rec.String()

	expected := fmt.Sprintf(`email: "%s"`, email)
	if strings.Index(str, expected) == -1 {
		t.Error("failed to found an information about an email")
	}

	expected = fmt.Sprintf(`text: "%s`, text)
	if strings.Index(str, expected) == -1 {
		t.Error("failed to found an information about a short text")
	}

	longText := "test text 01234567890123456789"
	rec = NewRecord(email, longText, nil)

	expected = fmt.Sprintf(`text: "%s`, longText[:8])
	if strings.Index(rec.String(), expected) == -1 {
		t.Error("failed to found an information about a long text")
	}
}

func TestRecord_NextIntervalEmpty(t *testing.T) {
	// EMPTY schedule
	rec := NewRecord("test@test.test", "test", nil)

	if str := rec.TryingString(); str != "" {
		t.Error("got trying string '", str, "', but should be empty")
	}

	duration, ok := rec.NextDuration(time.Now())
	if ok {
		t.Error("the next interval shouldn't exist, but it is")
	}
	if duration != 0 {
		t.Error("the next interval isn't zero, but should")
	}
}

func TestRecord_NextInterval(t *testing.T) {
	// check schedule
	rec := NewRecord("test@test.test", "test", []time.Duration{
		0,
		4 * time.Second,
		25 * time.Second,
		-4 * time.Second,
	})

	if str := rec.TryingString(); str != "" {
		t.Error("got trying string '", str, "', but should be empty for rec with a schedule")
	}

	// the last duration is checking getting value after a schedule list finished
	expectedDurations := []int{0, 4, 25, 0, 0}

	existDurations := make([]int, 0, len(expectedDurations))
	ok := true
	duration := 0 * time.Second
	start := time.Now()

	for ok {
		duration, ok = rec.NextDuration(start)

		v := float64(duration) / float64(time.Second)

		existDurations = append(existDurations, int(math.Ceil(v)))
	}

	if !reflect.DeepEqual(existDurations, expectedDurations) {
		t.Error("failed to get durations list for rec. Got ", existDurations, ", but expected is ", expectedDurations)
	}

	if str := rec.TryingString(); str == "" {
		t.Error("got an empty trying string, but should be with trying counts")
	}
}
