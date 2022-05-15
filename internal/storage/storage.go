package storage

import (
	"errors"
	"sync"
)

type Storage interface {
	GetGauge(metricName string) (float64, error)
	SetGauge(metricName string, metricValue float64)
	GetAllGauges() map[string]float64

	GetCounter(metricName string) (int64, error)
	SetCounter(metricName string, metricValue int64)
	GetAllCounters() map[string]int64
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
