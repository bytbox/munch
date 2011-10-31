package main

import (
	"encoding/hex"
	"log"
	"xml"
)

type RSSData struct {
	Channel Channel
}

type Channel struct {
	Title       string
	Link        string
	Description string
	Item        []RSSItemData
}

type RSSItemData struct {
	Title       string
	Link        string
	PubDate     string
	GUID        string
	Description string
	Content     string
}

func readRSS(feed Feed) {
	url := feed.Info.URL
	log.Print(url)
	r, err := client.Get(url)
	if err != nil {
		log.Print("ERROR: ", err.String())
		return
	}
	reader := r.Body
	feedData := RSSData{}
	err = xml.Unmarshal(reader, &feedData)
	if err != nil {
		log.Print("ERROR: ", err.String())
		return
	}
	// now transform the XML into our internal data structure
	changed := false
	for _, itemData := range feedData.Channel.Item {
		guid := itemData.GUID
		id := hex.EncodeToString([]byte(feed.Info.Name + guid))
		_, ok := feed.Items[guid]
		if !ok {
			// GUID not found - add the item
			changed = true
			t := parseTime(itemData.PubDate)
			item := Item {
				Feed: &feed,
				Title: itemData.Title,
				GUID: guid,
				URL: itemData.Link,
				Date: *t,
				Desc: itemData.Description,
				Content: itemData.Content,
				Read: false,
				ID: id,
			}
			feed.Items[guid] = item
		}
	}
	if (changed) {
		updates <- DoUpdate
		// TODO run some commands from Config?
	}
}


