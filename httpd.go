package main

import (
	"log"
	"net/http"
	"path"
)

func RunHTTPServer() {
	log.Print("Spawning HTTP Server")
	http.HandleFunc("/", HTTPHandler)
	http.HandleFunc("/open/", OpenHandler)
	http.HandleFunc("/about", AboutHandler)
	err := http.ListenAndServe(httpsrv, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

func HTTPHandler(w http.ResponseWriter, req *http.Request) {
	updates <- GetContent
	w.Write([]byte(<-content))
}

func OpenHandler(w http.ResponseWriter, req *http.Request) {
	// process the action and then redirect to /

	url := req.URL
	p := url.Path
	dir, itemid := path.Split(p)
	fid := path.Base(dir)
	var feed *Feed
	var item *Item

	// Find the feed with matching id
	for _, f := range Feeds {
		if fid == f.ID {
			feed = f
			break
		}
	}

	// Find the item with matching id
	for _, i := range feed.Items {
		if itemid == i.ID {
			item = i
			break
		}
	}

	// Update this item as read
	item.Read = true

	// perform updates
	updates <- DoUpdate

	// redirect to /
	http.Redirect(w, req, "/", http.StatusFound)
}

func AboutHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(about_string))
}
