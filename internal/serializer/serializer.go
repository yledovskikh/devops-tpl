package serializer

import (
	"encoding/json"
	"log"

	//"github.com/yledovskikh/devops-tpl/internal/storage"
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

func DecodingJSONMetric(b io.Reader) Metric {

	var m Metric
	err := json.NewDecoder(b).Decode(&m)
	if err != nil {
		log.Println("Error invalid decode request")
		return Metric{}
	}
	return m
}

func DecodingStringMetric(metricType, metricName, metricValue string) Metric {

	switch metricType {
	case "gauge":
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return Metric{MType: "invalidmetrictype"}
		}
		m := Metric{ID: metricName, MType: metricType, Value: &value}
		return m
	case "counter":
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return Metric{MType: "invalidmetrictype"}
		}
		m := Metric{ID: metricName, MType: metricType, Delta: &value}
		return m
	}
	//err := fmt.Errorf("unknown metric type %s, expected (gauge|counter)", metricType)
	return Metric{}
}

////CHECK Methods
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
