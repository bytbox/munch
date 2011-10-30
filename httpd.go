package main

import (
	"http"
	"log"
)

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

