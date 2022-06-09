package serializer

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/yledovskikh/devops-tpl/internal/hash"
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

//func SerializeGauge(metricName string, metricValue float64, h string) storage.Metric {
//	return storage.Metric{ID: metricName, MType: "gauge", Value: &metricValue, Hash: h}
//}

//func SerializeCounter(metricName string, metricValue int64, h string) storage.Metric {
//
//	return storage.Metric{ID: metricName, MType: "counter", Delta: &metricValue, Hash: h}
//}

func SerializeResponse(msg string) storage.JSONResponse {
	return storage.JSONResponse{Message: msg}
}

func SerializeGaugeH(metricName string, metricValue float64, key string) storage.Metric {
	var h string
	if key != "" {
		data := fmt.Sprintf("%s:gauge:%f", metricName, metricValue)
		h = hash.SignData(key, data)
	}
	return storage.Metric{ID: metricName, MType: "gauge", Value: &metricValue, Hash: h}
}

func SerializeCounterH(metricName string, metricValue int64, key string) storage.Metric {
	var h string
	if key != "" {
		data := fmt.Sprintf("%s:counter:%d", metricName, metricValue)
		h = hash.SignData(key, data)
	}
	return storage.Metric{ID: metricName, MType: "counter", Delta: &metricValue, Hash: h}
}
