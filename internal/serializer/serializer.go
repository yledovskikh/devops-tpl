package serializer

import (
	"encoding/json"
	"log"
	//"github.com/yledovskikh/devops-tpl/internal/storage"
	"io"
)

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type JsonResponse struct {
	Message string `json:"message"` // значение метрики в случае передачи gauge
}

//type Metric1 struct {
//	ID    string  `json:"id"`              // имя метрики
//	MType string  `json:"type"`            // параметр, принимающий значение gauge или counter
//	Delta int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
//	Value float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
//}

//type Metrics []Metric1

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

func EncodingMetricGauge(id string, value float64) ([]byte, error) {
	return json.Marshal(Metric{ID: id, MType: "gauge", Value: &value})
}

func EncodingMetricCounter(id string, value int64) ([]byte, error) {
	return json.Marshal(Metric{ID: id, MType: "counter", Delta: &value})
}

func DecodingResponse(msg string) JsonResponse {
	//message := make(map[string]string)
	//message["message"] = msg
	//resp, err := json.Marshal(message)
	//return resp, err
	return JsonResponse{Message: msg}
}

//
//func DecodingStringMetric(metricType, metricName, metricValue string) Metric {
//
//	switch metricType {
//	case "gauge":
//		value, err := strconv.ParseFloat(metricValue, 64)
//		if err != nil {
//			return Metric{MType: "notimplemented"}
//		}
//		m := Metric{ID: metricName, MType: metricType, Value: &value}
//		return m
//	case "counter":
//		value, err := strconv.ParseInt(metricValue, 10, 64)
//		if err != nil {
//			return Metric{MType: "notimplemented"}
//		}
//		m := Metric{ID: metricName, MType: metricType, Delta: &value}
//		return m
//	}
//	//err := fmt.Errorf("unknown metric type %s, expected (gauge|counter)", metricType)
//	return Metric{}
//}

//func EncodeServerResponse(resp ServerResponse) ([]byte, error) {
//	j, err := json.Marshal(resp)
//	if err != nil {
//		return nil, fmt.Errorf("failed to encode server response: %s", err.Error())
//	}
//	return j, nil
//
//}
