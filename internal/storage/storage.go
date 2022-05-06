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
}

var (
	mutex             = &sync.RWMutex{}
	ErrBadRequest     = errors.New("invalid value")
	ErrNotFound       = errors.New("metric not found")
	ErrNotImplemented = errors.New("unknown metric type")
)

type MetricStore struct {
	counter map[string]int64
	gauge   map[string]float64
}

func NewMetricStore() *MetricStore {
	return &MetricStore{
		counter: make(map[string]int64),
		gauge:   make(map[string]float64),
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
		s.gauge[metricName] = vg
	case "counter":
		vg, err := strconv.ParseInt(metricValue, 10, 64)

		if err != nil {
			return ErrBadRequest
		}
		s.counter[metricName] += vg
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
		if val, ok := s.gauge[metricName]; ok {
			return fmt.Sprintf("%v", val), nil
		}
	case "counter":
		if val, ok := s.counter[metricName]; ok {
			return fmt.Sprintf("%v", val), nil
		}
	}
	return "", ErrNotFound
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
