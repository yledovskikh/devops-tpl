package storage

import (
	"errors"
)

type Storage interface {
	SetMetrics(counters *map[string]int64, gauges *map[string]float64) error

	GetGauge(metricName string) (float64, error)
	SetGauge(metricName string, metricValue float64) error
	GetAllGauges() map[string]float64

	GetCounter(metricName string) (int64, error)
	SetCounter(metricName string, metricValue int64) error
	GetAllCounters() map[string]int64
	PingDB() error
	Close()
}

var (
	ErrBadRequest     = errors.New("invalid value")
	ErrNotFound       = errors.New("metric not found")
	ErrNotImplemented = errors.New("unknown metric type")
)
