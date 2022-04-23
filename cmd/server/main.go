package main

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type RunTimeMetrics struct {
	counter map[string]int64
	gauge   map[string]float64
}

var rtm RunTimeMetrics

func (rtm *RunTimeMetrics) UpdateRTMetric(metricType string, metricName string, metricValue string) error {
	switch strings.ToLower(metricType) {
	case "gauge":
		vg, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return errors.New("incorrect metric value")
		}
		if rtm.gauge == nil {
			rtm.gauge = make(map[string]float64)
		}
		rtm.gauge[metricName] = vg
	case "counter":
		vg, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return errors.New("incorrect metric value")
		}
		if rtm.counter == nil {
			rtm.counter = make(map[string]int64)
		}
		rtm.counter[metricName] += vg
	default:
		return errors.New("incorrect type (expected gauge or counter)")
	}
	return nil
}

func updateMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")
	err := rtm.UpdateRTMetric(metricType, metricName, metricValue)
	if err != nil {
		//TODO do correct error
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(rtm.counter, rtm.gauge)
}

func getMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	var metricValue string
	switch strings.ToLower(metricType) {
	case "gauge":
		if val, ok := rtm.gauge[metricName]; ok {
			metricValue = strconv.FormatFloat(val, 'f', 10, 64)
		}
	case "counter":
		if val, ok := rtm.counter[metricName]; ok {
			metricValue = strconv.FormatInt(val, 10)
		}
	}
	fmt.Println(metricValue)
	if metricValue == "" {
		http.Error(w, "404 page not found", http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintln(w, metricValue)
	if err != nil {
		panic("error write data client")
	}
	//metricValue := rtm
}

func main() {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/update/{metricType:(counter|gauge)}/{metricName:\\w+}/{metricValue:(\\w+}", updateMetric)
	r.Get("/value/{metricType:(counter|gauge)}/{metricName:\\w+}", getMetric)

	log.Fatal(http.ListenAndServe(":8080", r))
}
