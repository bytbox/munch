package main

import (
	"time"
)

var formats = []string{
	time.RFC1123Z,
	time.RFC1123,
	time.RFC3339,
}

func parseTime(s string) (t time.Time) {
	success := false
	var err error
	for _, fmt := range formats {
		t, err = time.Parse(fmt, s)
		if err == nil {
			success = true
			break
		}
	}
	if !success {
		t = time.Time{}
	}
	return
}
