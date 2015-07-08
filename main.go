package main

import (
	"flag"
	"log"
	"net/http"
)

// +build !appengine

var bindAddr = flag.String("addr", ":8080", "http listen address")

func main() {

	flag.Parse()

	log.Printf("Listening on %v", *bindAddr)
	log.Fatal(http.ListenAndServe(*bindAddr, nil))
}
