package main

import (
	"fmt"
	"time"
)

type Record struct {
	Email string `json:"email"`
	Text  string `json:"text"`

	success     bool
	tryingCount int
	schedule    []time.Duration
}

func NewRecord(email, text string, schedule []time.Duration) *Record {
	return &Record{
		Email:    email,
		Text:     text,
		schedule: schedule,
	}
}

func (o *Record) String() string {
	text := o.Text
	if len(o.Text) >= 16 {
		text = o.Text[:16] + "..."
	}

	return fmt.Sprintf(`Record{email: "%s"; text: "%s"}`, o.Email, text)
}

func (o *Record) IsValid() bool {
	return o.Email != "" && o.Text != ""
}

func (o *Record) TryingString() string {
	if o.tryingCount <= 1 {
		return ""
	}

	return fmt.Sprint(", trying ", o.tryingCount, " time(s)")
}

func (o *Record) NextDuration(start time.Time) (time.Duration, bool) {
	if o.tryingCount >= len(o.schedule) {
		return 0, false
	}

	trying := o.tryingCount
	o.tryingCount++

	if trying == 0 {
		return time.Duration(o.schedule[trying]), true
	}

	duration := o.schedule[trying] - time.Now().Sub(start)
	if duration < 0 {
		return 0, true
	}

	return duration, true
}
