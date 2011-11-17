package main

import (
	"log"
)

func readRDF(feed *Feed) {
	url := feed.Info.URL
	log.Print(url)
	r, err := client.Get(url)
	if err != nil {
		log.Print("ERROR: ", err.Error())
		return
	}
	_ = r.Body
	defer r.Body.Close()

	changed := false

	if changed {
		updates <- DoUpdate
		// TODO
	}
}
