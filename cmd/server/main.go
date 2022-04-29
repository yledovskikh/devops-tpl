package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yledovskikh/devops-tpl/internal/handlers"
	"log"
	"net/http"
)

func main() {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/update/{metricType}/{metricName}/{metricValue}", handlers.PostMetric)
	r.Get("/value/{metricType}/{metricName}", handlers.GetMetric)

	log.Fatal(http.ListenAndServe(":8080", r))
}
