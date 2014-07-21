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
	w.WriteHeader(200)
}

func main() {
	http.Handle("/dump/", &dumper{})
	http.Handle("/tests/", http.FileServer(http.Dir(".")))
	http.Handle("/files/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.ListenAndServe(":8000", nil)
}
