package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/foo/bar", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.RequestURI()
		log.Println("Incoming request for", url)
		for k, v := range r.Header {
			log.Println("---> Header:", k, " = ", v)
		}
	})
	log.Fatal(http.ListenAndServe("0.0.0.0:8090", nil))
}
