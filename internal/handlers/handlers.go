package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/yledovskikh/devops-tpl/internal/serializer"
	"github.com/yledovskikh/devops-tpl/internal/storage"
)

type Server struct {
	storage storage.Storage
}

func New(storage storage.Storage) *Server {
	return &Server{
		storage: storage,
	}
}

func SaveStoreDecodeMetric(m serializer.Metric, s storage.Storage) error {

	switch strings.ToLower(m.MType) {
	case "gauge":
		s.SetGauge(m.ID, *m.Value)
		return nil
	case "counter":
		s.SetCounter(m.ID, *m.Delta)
		return nil
	}
	return storage.ErrNotImplemented
}

func (s *Server) UpdateJSONMetric(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	status := http.StatusOK
	msg := "Metric saved"
	resp := serializer.DecodingResponse(msg)

	m := serializer.DecodingJSONMetric(r.Body)
	err := SaveStoreDecodeMetric(m, s.storage)
	if err != nil {
		status = storageErrToStatus(err)
		resp = serializer.DecodingResponse(err.Error())
	}

	w.WriteHeader(status)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf(err.Error())
	}
}

func (s *Server) getStorageJSONMetric(m serializer.Metric) (serializer.Metric, error) {

	switch strings.ToLower(m.MType) {
	case "gauge":
		value, err := s.storage.GetGauge(m.ID)
		m.Value = &value
		if err != nil {
			return serializer.Metric{}, err
		}
	case "counter":
		value, err := s.storage.GetCounter(m.ID)
		m.Delta = &value
		if err != nil {
			return serializer.Metric{}, err
		}
	default:
		return serializer.Metric{}, storage.ErrNotImplemented
	}

	return m, nil
}

func (s *Server) GetJSONMetric(w http.ResponseWriter, r *http.Request) {
	m := serializer.DecodingJSONMetric(r.Body)

	resp, err := s.getStorageJSONMetric(m)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		status := storageErrToStatus(err)
		respErr := serializer.DecodingResponse(err.Error())
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(respErr)
		return
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf(err.Error())
	}

}

func (s *Server) saveStringMetric(metricType, metricName, metricValue string) error {

	switch metricType {
	case "gauge":
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return storage.ErrBadRequest
		}
		s.storage.SetGauge(metricName, value)
		return nil
	case "counter":
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return storage.ErrBadRequest
		}
		s.storage.SetCounter(metricName, value)
		return nil
	}
	return storage.ErrNotImplemented
}

func (s *Server) getStringMetric(metricType, metricName string) (string, error) {

	switch metricType {
	case "gauge":
		value, err := s.storage.GetGauge(metricName)
		return fmt.Sprintf("%v", value), err
	case "counter":
		value, err := s.storage.GetCounter(metricName)
		return strconv.Itoa(int(value)), err
	}
	return "", storage.ErrNotFound
}

func (s *Server) UpdateURLMetric(w http.ResponseWriter, r *http.Request) {

	metricType := strings.ToLower(chi.URLParam(r, "metricType"))
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")
	err := s.saveStringMetric(metricType, metricName, metricValue)
	if err != nil {
		status := storageErrToStatus(err)
		w.WriteHeader(status)
	}
}

func (s *Server) GetURLMetric(w http.ResponseWriter, r *http.Request) {

	metricType := strings.ToLower(chi.URLParam(r, "metricType"))
	metricName := chi.URLParam(r, "metricName")

	metricValue, err := s.getStringMetric(metricType, metricName)
	if err == nil {
		fmt.Fprint(w, metricValue)
		return
	}
	status := storageErrToStatus(err)
	w.WriteHeader(status)
}

func storageErrToStatus(err error) int {
	switch err {
	case storage.ErrBadRequest:
		return http.StatusBadRequest
	case storage.ErrNotFound:
		return http.StatusNotFound
	case storage.ErrNotImplemented:
		return http.StatusNotImplemented
	default:
		return http.StatusInternalServerError
	}
}
