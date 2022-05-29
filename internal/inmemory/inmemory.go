package inmemory

import (
	"log"
	"sync"

	"github.com/yledovskikh/devops-tpl/internal/storage"
)

type MetricStore struct {
	counters     map[string]int64
	gauges       map[string]float64
	countersLock sync.RWMutex
	gaugesLock   sync.RWMutex
}

func NewMetricStore() *MetricStore {
	return &MetricStore{
		counters: make(map[string]int64),
		gauges:   make(map[string]float64),
	}
}

func (s *MetricStore) SetGauge(metricName string, metricValue float64) {
	s.gaugesLock.Lock()
	defer s.gaugesLock.Unlock()
	s.gauges[metricName] = metricValue
	log.Printf("save metric gauge - %s:%v", metricName, metricValue)

}

func (s *MetricStore) GetGauge(metricName string) (float64, error) {
	s.gaugesLock.RLock()
	defer s.gaugesLock.RUnlock()

	if val, ok := s.gauges[metricName]; ok {
		return val, nil
	}
	return 0, storage.ErrNotFound
}

func (s *MetricStore) GetAllGauges() map[string]float64 {
	s.gaugesLock.RLock()
	defer s.gaugesLock.RUnlock()
	m := make(map[string]float64)
	for i, val := range s.gauges {
		m[i] = val
	}
	return m
}

func (s *MetricStore) SetCounter(metricName string, metricValue int64) {
	s.countersLock.Lock()
	defer s.countersLock.Unlock()
	s.counters[metricName] += metricValue

	log.Printf("save metric counter - %s:%d", metricName, metricValue)

}

func (s *MetricStore) GetCounter(metricName string) (int64, error) {
	s.countersLock.RLock()
	defer s.countersLock.RUnlock()

	if val, ok := s.counters[metricName]; ok {
		return val, nil
	}
	return 0, storage.ErrNotFound
}

func (s *MetricStore) GetAllCounters() map[string]int64 {
	s.countersLock.RLock()
	defer s.countersLock.RUnlock()
	m := make(map[string]int64)
	for i, val := range s.counters {
		m[i] = val
	}
	return m
}
