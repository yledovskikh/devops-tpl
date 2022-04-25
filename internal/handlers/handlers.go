package handlers

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/yledovskikh/devops-tpl/internal/storage"
	"net/http"
	"strconv"
	"strings"
)

func updateMetric(metricType string, metricName string, metricValue string) error {
	switch strings.ToLower(metricType) {
	case "gauge":
		vg, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return errors.New("incorrect metric value")
		}
		//if storage.RunTimeMetrics == nil {
		//	storage.gauge = make(map[string]float64)
		//}
		storage.Gauge[metricName] = vg
	case "counter":
		vg, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return errors.New("incorrect metric value")
		}
		//if storage.counter == nil {
		//	storage.counter = make(map[string]int64)
		//}
		storage.Counter[metricName] += vg
	default:
		return errors.New("incorrect type (expected gauge or counter)")
	}
	return nil
}

func PostMetric(w http.ResponseWriter, r *http.Request) {

	//if r.Header.Get("Content-Type") != "text/plain" {
	//	http.Error(w, "Content-Type text/plain is required!", http.StatusUnsupportedMediaType)
	//	return
	//}

	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")
	if strings.ToLower(metricType) != "gauge" && strings.ToLower(metricType) != "counter" {
		http.Error(w, "incorrect metric type", http.StatusNotImplemented)
		return
	}
	err := updateMetric(metricType, metricName, metricValue)
	if err != nil {
		//TODO do correct error
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func GetMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	var metricValue string
	switch strings.ToLower(metricType) {
	case "gauge":
		if val, ok := storage.Gauge[metricName]; ok {
			metricValue = fmt.Sprintf("%v", val)
		}
	case "counter":
		if val, ok := storage.Counter[metricName]; ok {
			metricValue = fmt.Sprintf("%v", val)
		}
	}
	fmt.Println(metricValue)
	if metricValue == "" {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}
	_, err := fmt.Fprintln(w, metricValue)
	if err != nil {
		fmt.Println(err.Error())
	}
}
