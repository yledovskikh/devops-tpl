package main

import (
	h "github.com/yledovskikh/devops-tpl/internal/handlers"
	"log"
	"net/http"
)

func main() {
	http.Handle("/update/", http.HandlerFunc(h.MetricsHandler))
	server := &http.Server{
		Addr: ":8080",
	}
	log.Fatal(server.ListenAndServe())
}
