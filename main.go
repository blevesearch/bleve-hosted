//  Copyright (c) 2014 Couchbase, Inc.
//
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

// recreate the sample index
//go:generate rm -rf indexes/test.bleve
//go:generate bleve create indexes/test.bleve -s goleveldb
//go:generate bleve index indexes/test.bleve a.json

// +build !appengine,!appenginevm

package main

import (
	"flag"
	"log"
	"net/http"
)

var bindAddr = flag.String("addr", ":8080", "http listen address")

func main() {

	flag.Parse()

	log.Printf("Listening on %v", *bindAddr)
	log.Fatal(http.ListenAndServe(*bindAddr, nil))
}
