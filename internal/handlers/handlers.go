package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error function UpdateJSONMetric - ioutil.ReadAll(r.Body) - " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	m := serializer.DecodingJSONMetric(bytes.NewReader(b))

	err = SaveStoreDecodeMetric(m, s.storage)
	w.Header().Set("Content-Type", "application/json")

	var status int
	var resp []byte
	if err != nil {
		status = storageErrToStatus(err)
		resp, err = serializer.EncodingResponse(err.Error())
	} else {
		status = http.StatusOK
		msg := "Metric saved"
		resp, err = serializer.EncodingResponse(msg)
	}
	if err != nil {
		log.Println("Error function UpdateJSONMetric - serializer.EncodingResponse - ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Write(resp)
	//err = json.NewEncoder(w).Encode(resp)
	//if err != nil {
	//	log.Printf("Error UpdateJSONMetric - json.NewEncoder(w).Encode(resp) - %s", err.Error())
	//}
}

func (s *Server) getStorageJSONMetric(m serializer.Metric) (serializer.Metric, error) {
	//b, err := ioutil.ReadAll(r.Body)
	//if err != nil {
	//	err = errors.New("Error function getStorageJSONMetric - ioutil.ReadAll(r.Body) - " + err.Error())
	//	return serializer.Metric{}, err
	//}

	//m := serializer.DecodingJSONMetric(bytes.NewReader(b))
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

	//response, err := json.Marshal(m)
	//if err != nil {
	//	log.Println(err.Error())
	//	err = errors.New("Error function getStorageJSONMetric - json.Marshal(m) - " + err.Error())
	//	return nil, err
	//}
	//return response, err

}

func (s *Server) GetJSONMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	m := serializer.DecodingJSONMetric(r.Body)

	resp, err := s.getStorageJSONMetric(m)
	if err != nil {
		status := storageErrToStatus(err)
		respErr, e := serializer.EncodingResponse(err.Error())
		w.WriteHeader(status)
		w.Write(respErr)
		if e != nil {
			log.Printf("Error GetJSONMetric - serializer.EncodingResponse(err.Error()) - %s", err.Error())
			w.WriteHeader(status)
			return
		}
	}
	//w.WriteHeader(status)
	//w.Write(resp)
	err = json.NewEncoder(w).Encode(resp)
	//if err != nil {
	//	log.Printf("Error GetJSONMetric - json.NewEncoder(w).Encode(resp) - %s", err.Error())
	//}

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
