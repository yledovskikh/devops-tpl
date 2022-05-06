package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/yledovskikh/devops-tpl/internal/storage"
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

func (s *Server) PostMetric(w http.ResponseWriter, r *http.Request) {

	//if r.Header.Get("Content-Type") != "text/plain" {
	//	http.Error(w, "Content-Type text/plain is required!", http.StatusUnsupportedMediaType)
	//	return
	//}

	metricType := strings.ToLower(chi.URLParam(r, "metricType"))
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")
	err := s.storage.Put(metricType, metricName, metricValue)

	if err == nil {
		//w.WriteHeader(http.StatusOK)
		return
	}

	status := storageErrToStatus(err)
	w.WriteHeader(status)
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
	case storage.BadRequest:
		return http.StatusBadRequest
	case storage.NotFound:
		return http.StatusNotFound
	case storage.NotImplemented:
		return http.StatusNotImplemented
	default:
		return http.StatusInternalServerError
	}
}
