package reportmetrics

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/yledovskikh/devops-tpl/internal/storage"
)

func send2server(url string, metrics *[]storage.Metric) error {
	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(metrics); err != nil {
		log.Error().Err(err).Msg("")
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, payloadBuf)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		return err
	}

	err = response.Body.Close()
	if err != nil {
		return err
	}
	log.Info().Msgf("metrics was sent - status code: %d", response.StatusCode)
	return nil
}

func Exec(ctx context.Context, wg *sync.WaitGroup, endpoint string, ch <-chan *[]storage.Metric) {

	url := endpoint + "/updates/"

	for {
		select {
		case metrics := <-ch:

			err := send2server(url, metrics)
			if err != nil {
				log.Error().Err(err).Msg("")
				continue
			}
		case <-ctx.Done():
			log.Info().Msg("reporting of metrics was stopped")
			wg.Done()
			return
		}
	}
}
