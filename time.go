package main

import (
	"os"
	"time"
)

var formats = []string{
	time.RFC1123Z,
	time.RFC1123,
}

func parseTime(s string) (t *time.Time) {
	success := false
	var err os.Error
	for _, fmt := range formats {
		t, err = time.Parse(fmt, s)
		if err == nil {
			success = true
			break
		}
	}
	if !success {
		t = &time.Time{}
	}
	return
}

