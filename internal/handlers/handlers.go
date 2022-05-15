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
	//ms := map[string]string{"metricType": metricType, "metricName": metricName, "metricValue": metricValue}
	m := serializer.DecodingStringMetric(metricType, metricName, metricValue)
	err := s.storage.SetMetric(m)
	if err != nil {
		status := storageErrToStatus(err)
		w.WriteHeader(status)
	}
}

func (s *Server) UpdateJSONMetric(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error function UpdateJSONMetric ioutil.ReadAll(r.Body)")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	m := serializer.DecodingJSONMetric(bytes.NewReader(b))
	err = s.storage.SetMetric(m)

	if err != nil {
		status := storageErrToStatus(err)
		w.WriteHeader(status)
	}
}

func (s *Server) GetJSONMetric(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error read body: %s", err.Error())
	}

	m := serializer.DecodingJSONMetric(bytes.NewReader(b))
	//if err != nil {
	//	log.Printf("Error Descoding body: %s", err.Error())
	//	return
	//}

	metric, err := s.storage.GetMetric(m)

	if err == nil {
		response, err := json.Marshal(metric)
		if err != nil {
			log.Println(err.Error())
			status := http.StatusInternalServerError
			w.WriteHeader(status)
			return
		}
		//log.Println(response)

		w.Write(response)
		return
	}
	status := storageErrToStatus(err)
	w.WriteHeader(status)

}

func (s *Server) GetURLMetric(w http.ResponseWriter, r *http.Request) {

	metricType := strings.ToLower(chi.URLParam(r, "metricType"))
	metricName := chi.URLParam(r, "metricName")
	m := serializer.Metric{ID: metricName, MType: metricType}
	metric, err := s.storage.GetMetric(m)
	if err == nil {
		switch m.MType {
		case "gauge":
			fmt.Fprint(w, *metric.Value)
		case "counter":
			fmt.Fprint(w, *metric.Delta)
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
