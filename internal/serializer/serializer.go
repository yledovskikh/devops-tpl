package serializer

import (
	"encoding/json"
	"io"
	"log"

	"github.com/yledovskikh/devops-tpl/internal/storage"
)

func DecodingJSONMetric(r io.Reader) (storage.Metric, error) {

	var m storage.Metric
	err := json.NewDecoder(r).Decode(&m)
	if err != nil {
		log.Println("Error invalid decode request")
		return storage.Metric{}, err
	}
	return m, nil
}

func DecodingJSONMetrics(r io.Reader) ([]storage.Metric, error) {

	var m []storage.Metric
	err := json.NewDecoder(r).Decode(&m)
	if err != nil {
		log.Println("Error invalid decode request")
		return []storage.Metric{}, err
	}
	return m, nil
}

func SerializeGauge(metricName string, metricValue float64, h string) storage.Metric {
	return storage.Metric{ID: metricName, MType: "gauge", Value: &metricValue, Hash: h}
}

func SerializeCounter(metricName string, metricValue int64, h string) storage.Metric {

	return storage.Metric{ID: metricName, MType: "counter", Delta: &metricValue, Hash: h}
}

func SerializeResponse(msg string) storage.JSONResponse {
	return storage.JSONResponse{Message: msg}
}
