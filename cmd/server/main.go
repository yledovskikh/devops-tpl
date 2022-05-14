package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yledovskikh/devops-tpl/internal/handlers"
	"github.com/yledovskikh/devops-tpl/internal/storage"
	"log"
	"net/http"
)

func main() {

	r := chi.NewRouter()
	s := storage.NewMetricStore()
	h := handlers.New(s)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/update", h.PostJsonMetric)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", h.PostMetric)
	r.Get("/value/{metricType}/{metricName}", h.GetMetric)

	log.Fatal(http.ListenAndServe(":8080", r))
}
