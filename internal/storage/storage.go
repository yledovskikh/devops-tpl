package storage

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

////type RunTimeMetrics struct {
////TODO rewrite to interface
//var Counter = make(map[string]int64)
//var Gauge = make(map[string]float64)

type Storage interface {
	Get(metricType string, metricName string) (string, error)
	Put(metricType string, metricName string, metricValue string) error
	GetAllGauges() map[string]float64
	GetAllCounters() map[string]int64
}

var (
	mutex             = &sync.RWMutex{}
	ErrBadRequest     = errors.New("invalid value")
	ErrNotFound       = errors.New("metric not found")
	ErrNotImplemented = errors.New("unknown metric type")
)

type MetricStore struct {
	counters map[string]int64
	gauges   map[string]float64
}

//type Metrics struct {
//	ID    string   `json:"id"`              // имя метрики
//	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
//	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
//	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
//}

func NewMetricStore() *MetricStore {
	return &MetricStore{
		counters: make(map[string]int64),
		gauges:   make(map[string]float64),
	}
}

func (s *MetricStore) Put(metricType string, metricName string, metricValue string) error {
	mutex.Lock()
	defer mutex.Unlock()
	switch strings.ToLower(metricType) {
	case "gauge":
		vg, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return ErrBadRequest
		}
		s.gauges[metricName] = vg
	case "counter":
		vg, err := strconv.ParseInt(metricValue, 10, 64)

		if err != nil {
			return ErrBadRequest
		}
		s.counters[metricName] += vg
	default:
		return ErrNotImplemented
	}

	return nil
}

func (s *MetricStore) Get(metricType string, metricName string) (string, error) {

	mutex.RLock()
	defer mutex.RUnlock()

	switch metricType {
	case "gauge":
		if val, ok := s.gauges[metricName]; ok {
			return fmt.Sprintf("%v", val), nil
		}
	case "counter":
		if val, ok := s.counters[metricName]; ok {
			return fmt.Sprintf("%v", val), nil
		}
	}
	return "", ErrNotFound
}

func (s *MetricStore) GetAllGauges() map[string]float64 {
	mutex.RLock()
	g := make(map[string]float64)
	defer mutex.RUnlock()
	for k, v := range s.gauges {
		g[k] = v
	}
	return g
}

func (s *MetricStore) GetAllCounters() map[string]int64 {
	mutex.RLock()
	c := make(map[string]int64)
	defer mutex.RUnlock()
	for k, v := range s.counters {
		c[k] = v
	}
	return c
}

//}

//type gauge float64
//type counter int64

//type RunTimeMetrics struct {
//	Alloc         gauge
//	BuckHashSys   gauge
//	Frees         gauge
//	GCCPUFraction gauge
//	GCSys         gauge
//	HeapAlloc     gauge
//	HeapIdle      gauge
//	HeapInuse     gauge
//	HeapObjects   gauge
//	HeapReleased  gauge
//	HeapSys       gauge
//	LastGC        gauge
//	Lookups       gauge
//	MCacheInuse   gauge
//	MCacheSys     gauge
//	MSpanInuse    gauge
//	MSpanSys      gauge
//	Mallocs       gauge
//	NextGC        gauge
//	NumForcedGC   gauge
//	NumGC         gauge
//	OtherSys      gauge
//	PauseTotalNs  gauge
//	StackInuse    gauge
//	StackSys      gauge
//	Sys           gauge
//	TotalAlloc    gauge
//	RandomValue   gauge
//	PollCount     counter
//}
