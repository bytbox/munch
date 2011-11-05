package main

import (
	"log"
)

func readAtom(feed *Feed) {
	url := feed.Info.URL
	log.Print(url)
	r, err := client.Get(url)
	if err != nil {
		log.Print("ERROR: ", err.Error())
	}
	_ = r.Body

	changed := false

	if (changed) {
		updates <- DoUpdate
		// TODO
	}
}
