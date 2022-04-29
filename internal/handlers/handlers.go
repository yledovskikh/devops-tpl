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

func updateMetric(metricType string, metricName string, metricValue string) (int, error) {
	switch strings.ToLower(metricType) {
	case "gauge":
		vg, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return http.StatusBadRequest, fmt.Errorf("failed to convert %s to float64 - %s", metricValue, err.Error())
		}
		//TODO rewrite save metrics to "concurrency safe"
		storage.Gauge[metricName] = vg
	case "counter":
		vg, err := strconv.ParseInt(metricValue, 10, 64)

		if err != nil {
			return http.StatusBadRequest, fmt.Errorf("failed to convert %s to int64 - %s", metricValue, err.Error())
		}

		//TODO rewrite save metrics to "concurrency safe"
		storage.Counter[metricName] += vg
	default:
		return http.StatusNotImplemented, errors.New("unknown metric type (expected gauge or counter)")
	}
	return http.StatusOK, nil
}

func PostMetric(w http.ResponseWriter, r *http.Request) {

	//if r.Header.Get("Content-Type") != "text/plain" {
	//	http.Error(w, "Content-Type text/plain is required!", http.StatusUnsupportedMediaType)
	//	return
	//}

	metricType := strings.ToLower(chi.URLParam(r, "metricType"))
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")
	statusCode, err := updateMetric(metricType, metricName, metricValue)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}
	fmt.Fprintf(w, "%s/%s/%s - saved", metricType, metricName, metricValue)
}

func GetMetric(w http.ResponseWriter, r *http.Request) {

	metricType := strings.ToLower(chi.URLParam(r, "metricType"))
	metricName := chi.URLParam(r, "metricName")
	var metricValue string
	//TODO refactor - move check metricType to storage method
	switch metricType {
	case "gauge":
		if val, ok := storage.Gauge[metricName]; ok {
			metricValue = fmt.Sprintf("%v", val)
		}
	case "counter":
		if val, ok := storage.Counter[metricName]; ok {
			metricValue = fmt.Sprintf("%v", val)
		}
	default:
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
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
