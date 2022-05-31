package serializer

import (
	"encoding/json"
	"io"
	"log"
)

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
}

type JSONResponse struct {
	Message string `json:"message"` // значение метрики в случае передачи gauge
}

func DecodingJSONMetric(r io.Reader) (Metric, error) {

	var m Metric
	err := json.NewDecoder(r).Decode(&m)
	if err != nil {
		log.Println("Error invalid decode request")
		return Metric{}, err
	}
	return m, nil
}

func DecodingJSONMetrics(r io.Reader) ([]Metric, error) {

	var m []Metric
	err := json.NewDecoder(r).Decode(&m)
	if err != nil {
		log.Println("Error invalid decode request")
		return []Metric{}, err
	}
	return m, nil
}

func SerializeGauge(metricName string, metricValue float64, h string) Metric {
	return Metric{ID: metricName, MType: "gauge", Value: &metricValue, Hash: h}
}

func SerializeCounter(metricName string, metricValue int64, h string) Metric {

	return Metric{ID: metricName, MType: "counter", Delta: &metricValue, Hash: h}
}

func SerializeResponse(msg string) JSONResponse {
	return JSONResponse{Message: msg}
}
