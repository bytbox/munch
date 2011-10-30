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
	"template"
	"time"
	"xml"
)

type Configuration struct {
	UpdateInterval int64
	Feeds          []FeedInfo
}

type FeedInfo struct {
	Name string
	URL  string
	Type string
}

type Feed struct {
	Info  FeedInfo
	ID    string
	Items map[string]Item
}

type Item struct {
	Title   string
	GUID    string
	URL     string
	Date    string // TODO this may need to be a struct
	Desc    string
	Content string
	Read    bool
	ID	string
}

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

type TemplateData struct {
	Config *Configuration
	Feeds  map[string]Feed
}

const (
	DoUpdate = iota
	GetContent
)

var (
	Config        Configuration
	Feeds         map[string]Feed

	page_template template.Template
	client        http.Client
	updates       chan int
	content       chan string
)

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
		log.Fatal(err.String())
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatal(err.String())
	}
}

func InitCache() {
	Feeds = make(map[string]Feed)
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
			feed := Feed{}
			feed.Info = info
			feed.ID = hex.EncodeToString([]byte(name))
			feed.Items = make(map[string]Item)
			Feeds[name] = feed
		}
	}
}

func HandleUpdates() {
	pageBuffer := new(bytes.Buffer)
	tmplData := TemplateData{}
	tmplData.Feeds = Feeds
	tmplData.Config = &Config
	err := page_template.Execute(pageBuffer, tmplData)
	if err != nil {
		log.Print("ERROR: ", err.String())
	}
	for u := range updates {
		switch u {
		case DoUpdate:
			pageBuffer = new(bytes.Buffer)
			err = page_template.Execute(pageBuffer, tmplData)
			if err != nil {
				log.Print("ERROR: ", err.String())
			}
		case GetContent:
			content <- pageBuffer.String()
		default:
			panic("Undefined request to updater")
		}
	}
}

func WriteCache() {
	// TODO
}

type FeedReader interface {

}

func ReadFeeds() {
	log.Print("Updating feeds")
	for _, feed := range Feeds {
		switch (feed.Info.Type) {
		case "RSS":
			readRSS(feed)
		case "Atom":
		case "RDF":
		default:
			log.Print("Ignoring unknown feed of type ", feed.Info.Type)
		}
	}
	log.Print("Done")
}

// TODO handle this gracefully
func readFeed(feed Feed, reader FeedReader) {

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
			item := Item {
				Title: itemData.Title,
				GUID: guid,
				URL: itemData.Link,
				Date: itemData.PubDate,
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

func readAtom(feed Feed) {
	url := feed.Info.URL
	log.Print(url)
	r, err := client.Get(url)
	if err != nil {
		log.Print("ERROR: ", err.String())
	}
	_ = r.Body

	changed := false

	if (changed) {
		updates <- DoUpdate
		// TODO
	}
}

func readRDF(feed Feed) {
	url := feed.Info.URL
	log.Print(url)
	r, err := client.Get(url)
	if err != nil {
		log.Print("ERROR: ", err.String())
	}
	_ = r.Body

	changed := false

	if (changed) {
		updates <- DoUpdate
		// TODO
	}
}

func RunHTTPServer() {
	log.Print("Spawning HTTP Server")
	http.HandleFunc("/", HTTPHandler)
	err := http.ListenAndServe("localhost:8090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.String())
	}
}

func HTTPHandler(w http.ResponseWriter, req *http.Request) {
	updates <- GetContent
	w.Write([]byte(<-content))
}

