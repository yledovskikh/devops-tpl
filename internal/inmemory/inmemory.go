package inmemory

import (
	"errors"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/yledovskikh/devops-tpl/internal/storage"
)

type MetricStore struct {
	counters     map[string]int64
	gauges       map[string]float64
	countersLock sync.RWMutex
	gaugesLock   sync.RWMutex
	storage.Storage
}

func NewMetricStore() *MetricStore {
	return &MetricStore{
		counters: make(map[string]int64),
		gauges:   make(map[string]float64),
	}
}

func (s *MetricStore) SetGauge(metricName string, metricValue float64) error {
	s.gaugesLock.Lock()
	defer s.gaugesLock.Unlock()
	s.gauges[metricName] = metricValue
	log.Debug().Msgf("metric was saved metricType: gauges, metricName:%s, metricValue:%f", metricName, metricValue)

	return nil

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

func (s *MetricStore) SetCounter(metricName string, metricValue int64) error {
	s.countersLock.Lock()
	defer s.countersLock.Unlock()
	s.counters[metricName] += metricValue
	log.Debug().Msgf("metric was saved metricType: counters, metricName:%s, metricValue:%f", metricName, metricValue)
	return nil
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

func (s *MetricStore) PingDB() error {
	err := errors.New("not configured db")
	return err
}

func (s *MetricStore) Close() {
	//blanc func
}

func (s *MetricStore) SetMetrics(metrics *[]storage.Metric) error {

	for _, metric := range *metrics {
		switch strings.ToLower(metric.MType) {
		case "gauge":
			err := s.SetGauge(metric.ID, *metric.Value)
			if err != nil {
				log.Error().Err(err).Msg("")
			}
		case "counter":
			err := s.SetCounter(metric.ID, *metric.Delta)
			if err != nil {
				log.Error().Err(err).Msg("")
			}

		default:
			return storage.ErrNotFound
		}
	}
	return nil
}
