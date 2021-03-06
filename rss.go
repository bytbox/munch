package main

import (
	"encoding/hex"
	"encoding/xml"
	"log"
)

type RSSData struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Item        []RSSItemData `xml:"item"`
}

type RSSItemData struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
	Description string `xml:"description"`
	Content     string `xml:"content"`
}

func readRSS(feed *Feed) {
	url := feed.Info.URL
	r, err := client.Get(url)
	if err != nil {
		log.Print("ERROR fetching ", url, ": ", err.Error())
		return
	}
	reader := r.Body
	defer r.Body.Close()
	feedData := RSSData{}
	err = xml.Unmarshal(reader, &feedData)
	if err != nil {
		log.Print("ERROR parsing ", url, ": ", err.Error())
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
			item := &Item{
				Feed:    feed,
				Title:   itemData.Title,
				GUID:    guid,
				URL:     itemData.Link,
				Date:    t,
				Desc:    itemData.Description,
				Content: itemData.Content,
				Read:    false,
				ID:      id,
			}
			feed.Items[guid] = item
		}
	}
	if changed {
		updates <- DoUpdate
		// TODO run some commands from Config?
	}
}
