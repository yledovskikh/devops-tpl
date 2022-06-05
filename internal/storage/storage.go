package storage

import (
	"errors"
)

type Storage interface {
	SetMetrics(*[]Metric) error

	GetGauge(metricName string) (float64, error)
	SetGauge(metricName string, metricValue float64) error
	GetAllGauges() map[string]float64

	GetCounter(metricName string) (int64, error)
	SetCounter(metricName string, metricValue int64) error
	GetAllCounters() map[string]int64
	PingDB() error
	Close()
}

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
}

type JSONResponse struct {
	Message string `json:"message"` // значение метрики в случае передачи gauge
}

var (
	ErrBadRequest     = errors.New("invalid value")
	ErrNotFound       = errors.New("metric not found")
	ErrNotImplemented = errors.New("unknown metric type")
)
