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
func (s *Server) UpdateMetric(w http.ResponseWriter, r *http.Request) {

	metricType := strings.ToLower(chi.URLParam(r, "metricType"))
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")
	err := s.storage.Put(metricType, metricName, metricValue)

	if err == nil {
		return
	}

	status := storageErrToStatus(err)
	w.WriteHeader(status)
}

func (s *Server) UpdateJSONMetric(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//TODO Обработать корректноый статус
		w.WriteHeader(http.StatusBadRequest)
	}
	m, err := serializer.DecodingMetric(bytes.NewReader(b))
	if err != nil {
		//TODO Обработать корректноый статус
		w.WriteHeader(http.StatusBadRequest)
	}
	switch strings.ToLower(m.MType) {
	case "gauge":
		s.storage.PutGauge(m.ID, *m.Value)
		//TODO remove debug message
		value, _ := s.storage.Get(m.MType, m.ID)
		log.Printf("Debug gauge: \n metric name: %s value: %s \n", m.MType, value)

	case "counter":
		s.storage.PutCounter(m.ID, *m.Delta)
		//TODO remove debug message
		value, _ := s.storage.Get(m.MType, m.ID)
		log.Printf("Debug counter: \n metric name: %s value: %s \n", m.MType, value)
	}
}

func (s *Server) GetJSONMetric(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error read body: %s", err.Error())
	}

	m, err := serializer.DecodingMetric(bytes.NewReader(b))
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

func (s *Server) GetMetric(w http.ResponseWriter, r *http.Request) {

	metricType := strings.ToLower(chi.URLParam(r, "metricType"))
	metricName := chi.URLParam(r, "metricName")
	val, err := s.storage.Get(metricType, metricName)

	if err == nil {
		fmt.Fprint(w, val)
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
