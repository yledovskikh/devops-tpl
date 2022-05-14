package serializer

import (
	"encoding/json"
	"io"
)

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

//TODO Проверить нужно или нет
type ServerResponse struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

type Metrics []Metric

func EncodingMetricGauge(id string, value float64) ([]byte, error) {
	return json.Marshal(Metric{ID: id, MType: "gauge", Value: &value})
}

func EncodingMetricCounter(id string, value int64) ([]byte, error) {
	return json.Marshal(Metric{ID: id, MType: "counter", Delta: &value})
}

func DecodingMetric(b io.Reader) (Metric, error) {

	var m Metric
	err := json.NewDecoder(b).Decode(&m)
	if err != nil {
		return Metric{}, err
	}
	return m, err
}

//func EncodeServerResponse(resp ServerResponse) ([]byte, error) {
//	j, err := json.Marshal(resp)
//	if err != nil {
//		return nil, fmt.Errorf("failed to encode server response: %s", err.Error())
//	}
//	return j, nil
//
//}
