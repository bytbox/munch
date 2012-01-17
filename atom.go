package main

import (
	"encoding/hex"
	"encoding/xml"
	"log"
)

type AtomFeed struct {
	Title string `xml:"title"`
	Link  AtomLink `xml:"link"`
	Entry []AtomItemData `xml:"entity"`
}

type AtomItemData struct {
	ID      string `xml:"id"`
	Link    AtomLink `xml:"link"`
	Title   string `xml:"title"`
	Updated string `xml:"updated"`
	Summary string `xml:"summary"`
}

type AtomLink struct {
	Href string `xml:"attr"`
}

func readAtom(feed *Feed) {
	url := feed.Info.URL
	r, err := client.Get(url)
	if err != nil {
		log.Print("ERROR fetching ", url, ": ", err.Error())
		return
	}
	reader := r.Body
	defer r.Body.Close()
	feedData := AtomFeed{}
	err = xml.Unmarshal(reader, &feedData)
	if err != nil {
		log.Print("ERROR parsing ", url, ": ", err.Error())
		return
	}

	changed := false
	for _, itemData := range feedData.Entry {
		uid := itemData.ID
		id := hex.EncodeToString([]byte(feed.Info.Name + uid))
		_, ok := feed.Items[uid]
		if !ok {
			changed = true
			t := parseTime(itemData.Updated)
			item := &Item{
				Feed:  feed,
				Title: itemData.Title,
				GUID:  uid,
				URL:   itemData.Link.Href,
				Date:  t,
				Desc:  itemData.Summary,
				Read:  false,
				ID:    id,
			}
			feed.Items[uid] = item
		}
	}

	if changed {
		updates <- DoUpdate
		// TODO run some commands from Config?
	}
}
