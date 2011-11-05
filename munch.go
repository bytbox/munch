package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"http"
	"json"
	"log"
	"os"
	"path"
	"sort"
	"template"
	"time"
	"url"
)

type Configuration struct {
	UpdateInterval int64
	Feeds          []FeedInfo
}

type FeedInfo struct {
	Name    string
	URL     string
	MainURL string
	Type    string
}

type Feed struct {
	Info  FeedInfo
	ID    string
	Items map[string]*Item
}

type Item struct {
	Feed	*Feed
	Title   string
	GUID    string
	URL     string
	Date    time.Time
	Desc    string
	Content string
	Read    bool
	ID	string
}

type TemplateData struct {
	Config *Configuration
	Feeds  *map[string]*Feed
	Unread []*Item
}

const (
	DoUpdate = iota
	GetContent
)

var (
	Config        Configuration
	Feeds         map[string]*Feed

	page_template template.Template
	client        http.Client
	updates       chan int
	content       chan string
)

type FeedReader interface {

}

func main() {
	flag.Parse()

	updates = make(chan int)
	content = make(chan string)

	InitTemplate()
	ReadConfig()
	InitCache()
	go HandleUpdates()
	WriteCache()
	go RunHTTPServer()

	ticks := time.Tick(1e9 * Config.UpdateInterval)
	for {
		go ReadFeeds()
		<-ticks
	}
}

func InitTemplate() {
	log.Print("Initializing Page Template")
	page_template.Parse(page_template_string)
}

func ReadConfig() {
	log.Print("Reading Config")
	// Read config from ~/.munchrc
	file, err := os.Open(path.Join(os.Getenv("HOME"),".munchrc"))
	if err != nil {
		log.Fatal(err.Error())
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func InitCache() {
	Feeds = make(map[string]*Feed)
	// Ensure the cache directory exists
	cachePath := path.Join(os.Getenv("HOME"), ".munch.d", "cache")
	os.MkdirAll(cachePath, 0700)
	// For each feed
	for _, info := range Config.Feeds {
		name := info.Name
		fPath := path.Join(cachePath, name)
		file, _ := os.Open(fPath)
		if file != nil {
		} else {
			log.Print("New Feed: ", name)
			feed := &Feed{}
			feed.Info = info
			feed.ID = hex.EncodeToString([]byte(name))
			feed.Items = make(map[string]*Item)
			if feed.Info.MainURL == "" {
				u, err := url.Parse(feed.Info.URL)
				if err != nil {
					panic(err)
				}
				feed.Info.MainURL = "http://" + u.Host
			}
			Feeds[name] = feed
		}
	}
}

func HandleUpdates() {
	pageBuffer := new(bytes.Buffer)
	tmplData := TemplateData{}
	tmplData.Feeds = &Feeds
	tmplData.Config = &Config
	err := page_template.Execute(pageBuffer, tmplData)
	if err != nil {
		log.Print("ERROR: ", err.Error())
	}
	for u := range updates {
		switch u {
		case DoUpdate:
			tmplData.Unread = getUnread(tmplData.Feeds)
			pageBuffer = new(bytes.Buffer)
			err = page_template.Execute(pageBuffer, tmplData)
			if err != nil {
				log.Print("ERROR: ", err.Error())
			}
		case GetContent:
			content <- pageBuffer.String()
		default:
			panic("Undefined request to updater")
		}
	}
}

type ItemList []*Item

func (is ItemList) Len() int {
	return len(is)
}

func (is ItemList) Less(i, j int) bool {
	a := is[i]
	b := is[j]
	return a.Date.Seconds() > b.Date.Seconds()
}

func (is ItemList) Swap(i, j int) {
	is[i], is[j] = is[j], is[i]
}

func getUnread(feeds *map[string]*Feed) []*Item {
	items := ItemList{}
	for _, feed := range *feeds {
		for _, item := range feed.Items {
			if !item.Read {
				items = append(items, item)
			}
		}
	}
	sort.Sort(items)
	return items
}

func WriteCache() {
	// TODO
}

func ReadFeeds() {
	log.Print("Updating feeds")
	for _, feed := range Feeds {
		switch (feed.Info.Type) {
		case "RSS":
			readRSS(feed)
		case "Atom":
			readAtom(feed)
		case "RDF":
			readRDF(feed)
		default:
			log.Print("Ignoring unknown feed of type ", feed.Info.Type)
		}
	}
	log.Print("Done")
}

// TODO handle this gracefully
func readFeed(feed Feed, reader FeedReader) {

}

