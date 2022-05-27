package handlers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/yledovskikh/devops-tpl/internal/hash"
	"github.com/yledovskikh/devops-tpl/internal/serializer"
	"github.com/yledovskikh/devops-tpl/internal/storage"
)

type Server struct {
	storage storage.Storage
	Key     string
}

func New(storage storage.Storage) *Server {
	return &Server{
		storage: storage,
	}
}

func errJSONResponse(err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	status := storageErrToStatus(err)
	respErr := serializer.DecodingResponse(err.Error())
	w.WriteHeader(status)
	err = json.NewEncoder(w).Encode(respErr)
	if err != nil {
		log.Println(err)
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

func verifyMetric(m serializer.Metric, key string) error {
	var data string
	switch strings.ToLower(m.MType) {
	default:
		return storage.ErrNotImplemented
	case "gauge":
		data = fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value)
	case "counter":
		data = fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta)
	}
	h := hash.SignData(key, data)
	log.Println(h, data)
	v := hash.VerifyHash(m.Hash, h)
	if !v {
		return storage.ErrNotImplemented
	}
	return nil
}

func (s *Server) UpdateJSONMetric(w http.ResponseWriter, r *http.Request) {

	m, err := serializer.DecodingJSONMetric(r.Body)
	if err != nil {
		errJSONResponse(err, w)
		return
	}

	if s.Key != "" {
		err = verifyMetric(m, s.Key)
	}

	if err != nil {
		errJSONResponse(err, w)
		log.Println("Not verify hash")
		return
	}

	err = SaveStoreDecodeMetric(m, s.storage)
	if err != nil {
		errJSONResponse(err, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	msg := "Metric saved"
	resp := serializer.DecodingResponse(msg)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Print(err.Error())
	}
}

func (s *Server) getStorageJSONMetric(m serializer.Metric, key string) (serializer.Metric, error) {
	var data string
	switch strings.ToLower(m.MType) {
	case "gauge":
		value, err := s.storage.GetGauge(m.ID)
		m.Value = &value
		if err != nil {
			return serializer.Metric{}, err
		}
		data = fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value)
	case "counter":
		value, err := s.storage.GetCounter(m.ID)
		m.Delta = &value
		if err != nil {
			return serializer.Metric{}, err
		}
		data = fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta)
	default:
		return serializer.Metric{}, storage.ErrNotImplemented
	}
	if key != "" {
		h := hash.SignData(key, data)
		m.Hash = h
	}
	return m, nil
}

func (s *Server) GetJSONMetric(w http.ResponseWriter, r *http.Request) {

	m, err := serializer.DecodingJSONMetric(r.Body)
	if err != nil {
		errJSONResponse(err, w)
		return
	}
	resp, err := s.getStorageJSONMetric(m, s.Key)
	if err != nil {
		errJSONResponse(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Print(err.Error())
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
		_, err := fmt.Fprint(w, metricValue)
		if err != nil {
			log.Println(err)
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

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func CompressResponse(next http.Handler) http.Handler {
	// собираем Handler приведением типа
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// если gzip не поддерживается, передаём управление
			// дальше без изменений
			next.ServeHTTP(w, r)
			return
		}

		// создаём gzip.Writer поверх текущего w
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
		err = gz.Close()
		if err != nil {
			log.Println(err)
		}
	})
}

func DecompressRequest(next http.Handler) http.Handler {
	// собираем Handler приведением типа
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			// если gzip не поддерживается, передаём управление
			// дальше без изменений
			next.ServeHTTP(w, r)
			return
		}

		zr, err := gzip.NewReader(r.Body)
		if err != nil {
			err = fmt.Errorf("reading gzip body failed:%w", err)
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer zr.Close()
		r.Body = zr
		next.ServeHTTP(w, r)
	})
}

func (s *Server) AllMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	for metric, value := range s.storage.GetAllGauges() {
		fmt.Fprint(w, "<br>", metric, ":", value, "</br>")
	}
	for metric, value := range s.storage.GetAllCounters() {
		fmt.Fprint(w, "<br>", metric, ":", value, "</br>")
	}
}
