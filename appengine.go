//  Copyright (c) 2016 Couchbase, Inc.
//
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

// +build appengine

package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/blevesearch/bleve"
	bleveHttp "github.com/blevesearch/bleve/http"
	"github.com/blevesearch/bleve/index/store/gtreap"
	"github.com/blevesearch/blevex/preload"
	"github.com/gorilla/mux"
)

const indexDir = "indexes"

func init() {

	bleve.Config.DefaultKVStore = gtreap.Name

	bleveHttp.SetLog(log.New(os.Stderr, "bleve.http ", log.LstdFlags))

	router := mux.NewRouter()
	router.StrictSlash(true)

	listIndexesHandler := bleveHttp.NewListIndexesHandler()
	router.Handle("/api", listIndexesHandler).Methods("GET")

	docCountHandler := bleveHttp.NewDocCountHandler("")
	docCountHandler.IndexNameLookup = indexNameLookup
	router.Handle("/api/{indexName}/_count", docCountHandler).Methods("GET")

	searchHandler := bleveHttp.NewSearchHandler("")
	searchHandler.IndexNameLookup = indexNameLookup
	router.Handle("/api/{indexName}/_search", searchHandler).Methods("POST")

	http.Handle("/", &CORSWrapper{router})

	log.Printf("opening indexes")
	// walk the data dir and register index names
	dirEntries, err := ioutil.ReadDir(indexDir)
	if err != nil {
		log.Printf("error reading data dir: %v", err)
		return
	}

	for _, dirInfo := range dirEntries {
		indexPath := indexDir + string(os.PathSeparator) + dirInfo.Name()

		if dirInfo.IsDir() {
			log.Printf("see directory %s, this is not supported in the appengine environment", dirInfo.Name())
		} else {
			log.Printf("preloading index export %s", dirInfo.Name())
			// assume file in this dir is actually a bleve export
			i, err := bleve.NewUsing(
				"",
				bleve.NewIndexMapping(),
				bleve.Config.DefaultIndexType,
				preload.Name,
				map[string]interface{}{
					"kvStoreName_actual": gtreap.Name,
					"preloadpath":        indexPath,
				})
			if err != nil {
				log.Printf("error preloading index %s: %v", indexPath, err)
			} else {
				i.SetName(dirInfo.Name())
				log.Printf("registered index: %s", dirInfo.Name())
				bleveHttp.RegisterIndexName(dirInfo.Name(), i)
			}
		}
	}
}

func muxVariableLookup(req *http.Request, name string) string {
	return mux.Vars(req)[name]
}

func indexNameLookup(req *http.Request) string {
	return muxVariableLookup(req, "indexName")
}
