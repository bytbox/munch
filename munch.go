package main

import (
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
	Config Configuration
	Feeds  map[string]Feed
}

var (
	Config        Configuration
	Feeds         map[string]Feed
	TmplData      TemplateData

	page_template template.Template
	page_content  string
	client        http.Client
)

func main() {
	flag.Parse()

	InitTemplate()
	ReadConfig()
	InitCache()
	WriteCache()
	go RunHTTPServer()

	ticks := time.Tick(1e9 * Config.UpdateInterval)
	for {
		ReadFeeds()
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
			feed.Items = make(map[string]Item)
			Feeds[name] = feed
		}
	}
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
		default:
			log.Print("Ignoring unknown feed of type ", feed.Info.Type)
		}
	}
	log.Print("Done")
}

func readRSS(feed Feed) {
	url := feed.Info.URL
	log.Print(url)
	r, err := client.Get(url)
	if err != nil {
		log.Print("ERROR: ", err.String())
	}
	reader := r.Body
	feedData := RSSData{}
	err = xml.Unmarshal(reader, &feedData)
	if err != nil {
		log.Print("ERROR: ", err.String())
	}
	// now transform the XML into our internal data structure
	changed := false
	for _, itemData := range feedData.Channel.Item {
		guid := itemData.GUID
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
			}
			feed.Items[guid] = item
		}
	}
	if (changed) {
		UpdatePage()
		// TODO run some commands from Config?
	}
}

func UpdatePage() {
	// TODO create a write on page_content
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
	TmplData.Feeds = Feeds
	TmplData.Config = Config
	err := page_template.Execute(w, TmplData)
	if err != nil {
		log.Print("ERROR: ", err.String())
	}
}

