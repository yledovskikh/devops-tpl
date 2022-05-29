package storage

import (
	"errors"
)

type Storage interface {
	GetGauge(metricName string) (float64, error)
	SetGauge(metricName string, metricValue float64)
	GetAllGauges() map[string]float64

	GetCounter(metricName string) (int64, error)
	SetCounter(metricName string, metricValue int64)
	GetAllCounters() map[string]int64
}

var (
	ErrBadRequest     = errors.New("invalid value")
	ErrNotFound       = errors.New("metric not found")
	ErrNotImplemented = errors.New("unknown metric type")
)
