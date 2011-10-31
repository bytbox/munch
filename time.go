package main

import (
	"os"
	"time"
)

var formats = []string{
	"Mon, 02 Jan 2006 15:04:05 MST",
	"Mon, 02 Jan 2006 15:04:05 -0700",
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

