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
}

type JSONResponse struct {
	Message string `json:"message"` // значение метрики в случае передачи gauge
}

func DecodingJSONMetric(b io.Reader) Metric {

	var m Metric
	err := json.NewDecoder(b).Decode(&m)
	if err != nil {
		log.Println("Error invalid decode request")
		return Metric{}
	}
	return m
}

func DecodingGauge(metricName string, metricValue float64) Metric {
	return Metric{ID: metricName, MType: "gauge", Value: &metricValue}
}

func DecodingCounter(metricName string, metricValue int64) Metric {
	return Metric{ID: metricName, MType: "counter", Delta: &metricValue}
}

func DecodingResponse(msg string) JSONResponse {
	return JSONResponse{Message: msg}
}
