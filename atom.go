package main

import (
	"encoding/hex"
	"encoding/xml"
	"log"
)

type AtomFeed struct {
	Title string
	Link  AtomLink
	Entry []AtomItemData
}

type AtomItemData struct {
	ID      string
	Link    AtomLink
	Title   string
	Updated string
	Summary string
}

type AtomLink struct {
	Href string `xml:"attr"`
}

func readAtom(feed *Feed) {
	url := feed.Info.URL
	r, err := client.Get(url)
	if err != nil {
		log.Print("ERROR fetching ", url, ": ", err.Error())
	}
	reader := r.Body
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
				Date:  *t,
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
