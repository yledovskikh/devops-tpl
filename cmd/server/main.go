package main

import (
	"log"
	"net/http"
)

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only Post requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "text/plain; charset=UTF-8" {
		http.Error(w, r.Header.Get("Content-Type"), http.StatusUnsupportedMediaType)
		//http.Error(w, "Content-Type text/plain is required!", http.StatusUnsupportedMediaType)
		return
	}
	//fmt.Fprintln(w, r.URL)
}

func main() {
	http.Handle("/update/", http.HandlerFunc(metricsHandler))
	server := &http.Server{
		Addr: ":8080",
	}
	log.Fatal(server.ListenAndServe())
}
