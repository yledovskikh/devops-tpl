package handlers

import (
	"bytes"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/yledovskikh/devops-tpl/internal/serializer"
	"github.com/yledovskikh/devops-tpl/internal/storage"
	"io/ioutil"
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

//func UpdateMetricJSONHandler(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json; charset=utf-8")
//
//	requestCompressed :=
//		strings.Contains(strings.Join(r.Header["Content-Encoding"], ","), "gzip")
//	compressResponse :=
//		strings.Contains(strings.Join(r.Header["Accept-Encoding"], ","), "gzip")
//
//	b, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		e := fmt.Sprintf("failed to read body: %s", err.Error())
//		sendResponse(w, http.StatusBadRequest, serializer.ServerResponse{Error: e}, compressResponse)
//	}
//
//	if requestCompressed {
//		b, err = archive.Decompress(b)
//		if err != nil {
//			e := fmt.Sprintf("Failed to decompress request body: %s", err.Error())
//			sendResponse(w, http.StatusBadRequest, serializer.ServerResponse{Error: e}, compressResponse)
//			return
//		}
//	}
//
//	m, err := serializer.DecodeBody(bytes.NewReader(b))
//	if err != nil {
//		e := fmt.Sprintf("Failed to decode request body: %s", err.Error())
//		sendResponse(w, http.StatusBadRequest, serializer.ServerResponse{Error: e}, compressResponse)
//		return
//	}
//	updateMetric(m)
//	sendResponse(w, http.StatusOK, serializer.ServerResponse{Result: "metric was saved"}, compressResponse)
//
//}

func (s *Server) PostMetric(w http.ResponseWriter, r *http.Request) {

	//if r.Header.Get("Content-Type") != "text/plain" {
	//	http.Error(w, "Content-Type text/plain is required!", http.StatusUnsupportedMediaType)
	//	return
	//}
	//
	//m, err := serializer.DecodeBody(bytes.NewReader(b))
	////	if err != nil {
	////		e := fmt.Sprintf("Failed to decode request body: %s", err.Error())
	////		sendResponse(w, http.StatusBadRequest, serializer.ServerResponse{Error: e}, compressResponse)
	////		return
	////	}
	////	updateMetric(m)
	////	sendResponse(w, http.StatusOK, serializer.ServerResponse{Result: "metric was saved"}, compressResponse)

	//err := s.storage.Put()
	//
	//if err == nil {
	//	//w.WriteHeader(http.StatusOK)
	//	return
	//}
	//
	//status := storageErrToStatus(err)
	//w.WriteHeader(status)

	w.Header().Set("Content-Type", "application/json")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//TODO Обработать корректноый статус
		w.WriteHeader(http.StatusBadRequest)
	}
	m, err := serializer.DecodingMetric(bytes.NewReader(b))

	switch strings.ToLower(m.MType) {
	case "gauge":
		s.storage.PutGauge(m.ID, *m.Value)
		//TODO remove debug message
		m, _ := s.storage.Get(m.MType, m.ID)
		fmt.Println("Debug: \n", m)

	case "counter":
		s.storage.PutCounter(m.ID, *m.Delta)
		//TODO remove debug message
		m, _ := s.storage.Get(m.MType, m.ID)
		fmt.Println("Debug: \n", m)
	}
	return
}

//func sendResponse(w http.ResponseWriter, code int, resp serializer.ServerResponse) {
//	responseBody, err := serializer.EncodeServerResponse(resp)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		w.Write([]byte(fmt.Sprintf("failed to encode server response: %s", err.Error())))
//		return
//	}
//
//	w.WriteHeader(code)
//	w.Write(responseBody)
//}

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
