package reportmetrics

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/yledovskikh/devops-tpl/internal/storage"
)

func Exec(ctx context.Context, endpoint string, ch <-chan *[]storage.Metric) {

	url := endpoint + "/updates/"
	//var metrics *[]storage.Metric
	for {
		select {
		case metrics := <-ch:
			payloadBuf := new(bytes.Buffer)
			if err := json.NewEncoder(payloadBuf).Encode(metrics); err != nil {
				log.Error().Err(err).Msg("")
			}

			client := &http.Client{}
			req, err := http.NewRequest(http.MethodPost, url, payloadBuf)
			if err != nil {
				//log.Println(err)
				log.Error().Err(err).Msg("")
			}

			req.Header.Add("Content-Type", "application/json")
			response, err := client.Do(req)
			if err != nil {
				log.Error().Err(err).Msg("")
				continue
			}

			err = response.Body.Close()
			if err != nil {
				log.Error().Err(err).Msg("")
				continue
			}
			log.Info().Msgf("metrics was sent - status code: %d", response.StatusCode)
		case <-ctx.Done():
			log.Info().Msg("reporting of metrics stopped")
			return
		}
	}
}

//func reportMetrics(endpoint string, key string, ch <-chan *storage.Storage) {
//	url := endpoint + "/updates/"
//	var metrics []storage.Metric
//	for {
//		select {
//		case 	st:=<-ch:
//
//
//		}
//
//	}
//
//
//
//	for mName, mValue := range st.GetAllGauges() {
//		var h string
//		if key != "" {
//			data := fmt.Sprintf("%s:gauge:%f", mName, mValue)
//			h = hash.SignData(key, data)
//		}
//
//		m := serializer.SerializeGauge(mName, mValue, h)
//		metrics = append(metrics, m)
//	}
//
//	for mName, mValue := range a.storage.GetAllCounters() {
//		var h string
//		if key != "" {
//			data := fmt.Sprintf("%s:counter:%d", mName, mValue)
//			h = hash.SignData(key, data)
//		}
//		m := serializer.SerializeCounter(mName, mValue, h)
//		metrics = append(metrics, m)
//	}
//
//	payloadBuf := new(bytes.Buffer)
//	if err := json.NewEncoder(payloadBuf).Encode(metrics); err != nil {
//		log.Error().Err(err).Msg("")
//	}
//
//	if err := send2server(url, payloadBuf); err != nil {
//		log.Error().Err(err).Msg("")
//	}
//	//reset counter after send to server
//	err := a.storage.SetCounter("PollCount", 0)
//	if err != nil {
//		log.Error().Err(err).Msg("")
//	}
//}
