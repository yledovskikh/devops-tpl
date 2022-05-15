package storage

import (
	"errors"
	"github.com/yledovskikh/devops-tpl/internal/serializer"
	"log"
	"strings"
	"sync"
)

type Storage interface {
	GetGauge(metricName string) (float64, error)
	SetGauge(metricName string, metricValue float64)
	GetAllGauges() map[string]float64

	GetCounter(metricName string) (int64, error)
	SetCounter(metricName string, metricValue int64)
	GetAllCounters() map[string]int64

	SetMetric(m serializer.Metric) error
	GetMetric(m serializer.Metric) (serializer.Metric, error)
}

type MetricStore struct {
	counters map[string]int64
	gauges   map[string]float64
}

var (
	mutex             = &sync.RWMutex{}
	ErrBadRequest     = errors.New("invalid value")
	ErrNotFound       = errors.New("metric not found")
	ErrNotImplemented = errors.New("unknown metric type")
)

func NewMetricStore() *MetricStore {
	return &MetricStore{
		counters: make(map[string]int64),
		gauges:   make(map[string]float64),
	}
}

func (s *MetricStore) SetGauge(metricName string, metricValue float64) {
	mutex.Lock()
	defer mutex.Unlock()
	s.gauges[metricName] = metricValue
	log.Printf("save metric gauge - %s:%f", metricName, metricValue)

}

func (s *MetricStore) GetGauge(metricName string) (float64, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if val, ok := s.gauges[metricName]; ok {
		return val, nil
	}
	return 0, ErrNotFound
}

func (s *MetricStore) GetAllGauges() map[string]float64 {
	mutex.RLock()
	defer mutex.RUnlock()
	m := make(map[string]float64)
	for i, val := range s.gauges {
		m[i] = val
	}
	return m
}

func (s *MetricStore) SetCounter(metricName string, metricValue int64) {
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := s.counters[metricName]; ok {
		s.counters[metricName] += metricValue
	} else {
		s.counters[metricName] = metricValue
	}
	log.Printf("save metric counter - %s:%d", metricName, metricValue)

}

func (s *MetricStore) GetCounter(metricName string) (int64, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if val, ok := s.counters[metricName]; ok {
		return val, nil
	}
	return 0, ErrNotFound
}

func (s *MetricStore) GetAllCounters() map[string]int64 {
	mutex.RLock()
	defer mutex.RUnlock()
	m := make(map[string]int64)
	for i, val := range s.counters {
		m[i] = val
	}
	return m
}

func (s *MetricStore) SetMetric(m serializer.Metric) error {
	switch strings.ToLower(m.MType) {
	case "gauge":
		s.SetGauge(m.ID, *m.Value)
		return nil
	case "counter":
		s.SetCounter(m.ID, *m.Delta)
		return nil
	case "notimplemented":
		return ErrBadRequest
	}
	return ErrNotImplemented
}

func (s *MetricStore) GetMetric(m serializer.Metric) (serializer.Metric, error) {

	var metric serializer.Metric

	switch m.MType {
	case "gauge":
		val, err := s.GetGauge(m.ID)
		if err != nil {
			log.Printf("Error get metrics: %s, %s, %s", m.MType, m.ID, err.Error())
			return serializer.Metric{}, err
		}
		metric = serializer.Metric{ID: m.ID, MType: m.MType, Value: &val}
	case "counter":
		val, err := s.GetCounter(m.ID)
		if err != nil {
			log.Printf("Error get metrics: %s, %s, %s", m.MType, m.ID, err.Error())
			return serializer.Metric{}, err
		}
		metric = serializer.Metric{ID: m.ID, MType: m.MType, Delta: &val}
	default:
		err := ErrNotFound
		return serializer.Metric{}, err
	}
	return metric, nil

}

//func (s *MetricStore) SetMetric(metricType, metricName, metricValue string) {
//	switch strings.ToLower(metricType) {
//	case "gauge":
//
//		s.SetGauge(metricName, *m.Value)
//		log.Printf("save metric %s:%d", m.ID, m.Value)
//	case "counter":
//		s.storage.SetCounter(m.ID, *m.Delta)
//		log.Printf("save metric %s:%d", m.ID, m.Delta)
//	}
//}
