package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/yledovskikh/devops-tpl/internal/serializer"
	"github.com/yledovskikh/devops-tpl/internal/storage"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Server struct {
	storage storage.Storage
}

func New(storage storage.Storage) *Server {
	return &Server{
		storage: storage,
	}
}

func (s *Server) saveMetric(m serializer.Metric) {
	switch strings.ToLower(m.MType) {
	case "gauge":
		s.storage.SetGauge(m.ID, *m.Value)
		log.Printf("save metric %s:%d", m.ID, m.Value)
	case "counter":
		s.storage.SetCounter(m.ID, *m.Delta)
		log.Printf("save metric %s:%d", m.ID, m.Delta)
	}
}

func (s *Server) UpdateMetric(w http.ResponseWriter, r *http.Request) {

	metricType := strings.ToLower(chi.URLParam(r, "metricType"))
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")
	ms := map[string]string{"metricType": metricType, "metricName": metricName, "metricValue": metricValue}
	m, err := serializer.DecodingStringMetric(ms)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	s.saveMetric(m)
}

func (s *Server) UpdateJSONMetric(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	m, err := serializer.DecodingJSONMetric(bytes.NewReader(b))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	s.saveMetric(m)

}

func (s *Server) GetJSONMetric(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error read body: %s", err.Error())
	}

	m, err := serializer.DecodingJSONMetric(bytes.NewReader(b))
	if err != nil {
		log.Printf("Error Descoding body: %s", err.Error())
		return
	}

	var metric serializer.Metric

	switch m.MType {
	case "gauge":
		val, err := s.storage.GetGauge(m.ID)
		if err != nil {
			log.Printf("Error get metrics: %s, %s, %s", m.MType, m.ID, err.Error())
			return
		}
		metric = serializer.Metric{ID: m.ID, MType: m.MType, Value: &val}
	case "counter":
		val, err := s.storage.GetCounter(m.ID)
		if err != nil {
			log.Printf("Error get metrics: %s, %s, %s", m.MType, m.ID, err.Error())
			return
		}
		metric = serializer.Metric{ID: m.ID, MType: m.MType, Delta: &val}

	}

	response, err := json.Marshal(metric)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(response)

	_, err = w.Write(response)
	if err != nil {
		log.Printf("Error write client: %s", err.Error())
	}

}

func (s *Server) getMetric(m serializer.Metric) (serializer.Metric, error) {
	var metric serializer.Metric

	switch m.MType {
	case "gauge":
		val, err := s.storage.GetGauge(m.ID)
		if err != nil {
			log.Printf("Error get metrics: %s, %s, %s", m.MType, m.ID, err.Error())
			return serializer.Metric{}, err
		}
		metric = serializer.Metric{ID: m.ID, MType: m.MType, Value: &val}
	case "counter":
		val, err := s.storage.GetCounter(m.ID)
		if err != nil {
			log.Printf("Error get metrics: %s, %s, %s", m.MType, m.ID, err.Error())
			return serializer.Metric{}, err
		}
		metric = serializer.Metric{ID: m.ID, MType: m.MType, Delta: &val}
	}
	return metric, nil
}

func (s *Server) GetURLMetric(w http.ResponseWriter, r *http.Request) {

	metricType := strings.ToLower(chi.URLParam(r, "metricType"))
	metricName := chi.URLParam(r, "metricName")
	ms := map[string]string{"metricType": metricType, "metricName": metricName}
	m, err := serializer.DecodingStringMetric(ms)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	metric, err := s.getMetric(m)

	if err == nil {
		switch m.MType {
		case "gauge":
			fmt.Fprint(w, metric.Value)
		case "counter":
			fmt.Fprint(w, metric.Delta)
		}
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
