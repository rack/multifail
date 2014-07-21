/*
Copyright 2014 Google Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
)

type dumper struct{}

// ServeHTTP writes request dumps to a file named after the request URL and
// User-Agent header (URI escaped), or reads them. Reads also use the User-Agent
// header value for routing.
func (*dumper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ua := r.Header.Get("user-agent")
	if ua == "" {
		log.Print("no user agent provided")
		w.WriteHeader(400)
		return
	}

	experiment := filepath.Clean(r.URL.Path)
	if experiment == "." {
		log.Printf("no or bad experiment path given: %s", r.URL.Path)
		w.WriteHeader(400)
		return
	}

	fpath := filepath.Join(".", experiment, url.QueryEscape(ua))

	switch r.Method {
	case "POST":
	case "GET":
		http.ServeFile(w, r, fpath)
		return
	default:
		log.Print("request to dumper was not a POST")
		w.WriteHeader(500)
		return
	}

	if err := os.MkdirAll(filepath.Dir(fpath), 0777); err != nil {
		log.Printf("errror creating dir %s: %s", filepath.Dir(fpath), err)
		w.WriteHeader(500)
		return
	}
	f, err := os.Create(fpath)
	if err != nil {
		log.Printf("error creating %s: %s", fpath, err)
		w.WriteHeader(500)
		return
	}
	defer f.Close()

	b, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Printf("error dumping request: %s", err)
		w.WriteHeader(500)
		return
	}
	f.Write(b)

	log.Printf("wrote dump to %s", fpath)
	http.ServeFile(w, r, fpath)
}

func main() {
	http.Handle("/dump/", &dumper{})
	http.Handle("/tests/", http.FileServer(http.Dir(".")))
	http.Handle("/files/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
