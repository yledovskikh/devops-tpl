package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yledovskikh/devops-tpl/internal/config"
	"github.com/yledovskikh/devops-tpl/internal/handlers"
	"github.com/yledovskikh/devops-tpl/internal/storage"
	"log"
	"net/http"
)

func main() {
	endpoint := config.ServerConfig()
	r := chi.NewRouter()
	s := storage.NewMetricStore()
	h := handlers.New(s)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/update/", h.UpdateJSONMetric)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", h.UpdateMetric)
	r.Get("/value/{metricType}/{metricName}", h.GetURLMetric)
	r.Post("/value/", h.GetJSONMetric)

	log.Fatal(http.ListenAndServe(endpoint, r))
}
