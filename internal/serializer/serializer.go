package serializer

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type Metrics []Metric

func DecodingJSONMetric(b io.Reader) (Metric, error) {

	var m Metric
	err := json.NewDecoder(b).Decode(&m)
	if err != nil {
		return Metric{}, err
	}
	return m, err
}

func DecodingStringMetric(m map[string]string) (Metric, error) {

	switch m["metricType"] {
	case "gauge":
		value, err := strconv.ParseFloat(m["metricValue"], 64)
		if err != nil {
			return Metric{}, err
		}
		m := Metric{ID: m["metricName"], MType: m["metricType"], Value: &value}
		return m, nil
	case "counter":
		value, err := strconv.ParseInt(m["metricValue"], 10, 64)
		if err != nil {
			return Metric{}, err
		}
		m := Metric{ID: m["metricName"], MType: m["metricType"], Delta: &value}
		return m, nil
	}
	err := fmt.Errorf("unknown metric type %s, expected (gauge|counter)", m["metricType"])
	return Metric{}, err
}

func EncodingMetricGauge(id string, value float64) ([]byte, error) {
	return json.Marshal(Metric{ID: id, MType: "gauge", Value: &value})
}

func EncodingMetricCounter(id string, value int64) ([]byte, error) {
	return json.Marshal(Metric{ID: id, MType: "counter", Delta: &value})
}

//func EncodeServerResponse(resp ServerResponse) ([]byte, error) {
//	j, err := json.Marshal(resp)
//	if err != nil {
//		return nil, fmt.Errorf("failed to encode server response: %s", err.Error())
//	}
//	return j, nil
//
//}
