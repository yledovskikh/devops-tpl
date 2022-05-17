package agent

import (
	"bytes"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/yledovskikh/devops-tpl/internal/serializer"
	"github.com/yledovskikh/devops-tpl/internal/storage"
)

type Agent struct {
	storage storage.Storage
}

func New(storage storage.Storage) *Agent {
	return &Agent{
		storage: storage,
	}
}

func (a *Agent) collectMetrics() {

	rand.Seed(time.Now().UnixNano())
	r := rand.Float64()

	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	a.storage.SetGauge("Alloc", float64(rtm.Alloc))
	a.storage.SetGauge("BuckHashSys", float64(rtm.BuckHashSys))
	a.storage.SetGauge("Frees", float64(rtm.Frees))
	a.storage.SetGauge("GCCPUFraction", float64(rtm.GCCPUFraction))
	a.storage.SetGauge("GCSys", float64(rtm.GCSys))
	a.storage.SetGauge("HeapAlloc", float64(rtm.HeapAlloc))
	a.storage.SetGauge("HeapIdle", float64(rtm.HeapIdle))
	a.storage.SetGauge("HeapInuse", float64(rtm.HeapInuse))
	a.storage.SetGauge("HeapObjects", float64(rtm.HeapObjects))
	a.storage.SetGauge("HeapReleased", float64(rtm.HeapReleased))
	a.storage.SetGauge("HeapSys", float64(rtm.HeapSys))
	a.storage.SetGauge("LastGC", float64(rtm.LastGC))
	a.storage.SetGauge("Lookups", float64(rtm.Lookups))
	a.storage.SetGauge("MCacheInuse", float64(rtm.MCacheInuse))
	a.storage.SetGauge("MCacheSys", float64(rtm.MCacheSys))
	a.storage.SetGauge("MSpanSys", float64(rtm.MSpanSys))
	a.storage.SetGauge("Mallocs", float64(rtm.Mallocs))
	a.storage.SetGauge("NextGC", float64(rtm.NextGC))
	a.storage.SetGauge("NumForcedGC", float64(rtm.NumForcedGC))
	a.storage.SetGauge("NextGC", float64(rtm.NumGC))
	a.storage.SetGauge("OtherSys", float64(rtm.OtherSys))
	a.storage.SetGauge("PauseTotalNs", float64(rtm.PauseTotalNs))
	a.storage.SetGauge("StackInuse", float64(rtm.StackInuse))
	a.storage.SetGauge("StackSys", float64(rtm.StackSys))
	a.storage.SetGauge("Sys", float64(rtm.Sys))
	a.storage.SetGauge("TotalAlloc", float64(rtm.TotalAlloc))
	a.storage.SetGauge("MSpanInuse", float64(rtm.MSpanInuse))
	a.storage.SetGauge("NumGC", float64(rtm.NumGC))
	a.storage.SetGauge("RandomValue", r)
	a.storage.SetCounter("PollCount", 1)
	log.Println("INFO collect metrics")
}
func send2server(endpoint string, body []byte) error {

	url := endpoint + "/update/"
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
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
	log.Printf("INFO metric %s was sent to %s \n", body, url)
	return nil
}

func (a *Agent) postMetrics(endpoint string) {

	for mName, mValue := range a.storage.GetAllGauges() {
		body, err := serializer.EncodingMetricGauge(mName, mValue)
		if err != nil {
			log.Println(err)
		}
		if err = send2server(endpoint, body); err != nil {
			log.Println(err.Error())
			continue
		}
	}

	for mName, mValue := range a.storage.GetAllCounters() {
		body, err := serializer.EncodingMetricCounter(mName, mValue)

		if err != nil {
			log.Println(err.Error())
		}

		if err = send2server(endpoint, body); err != nil {
			log.Println(err)
			continue
		}
	}
}

func (a *Agent) Exec(endpoint string, pollInterval, reportInterval time.Duration) {
	pollIntervalTicker := time.NewTicker(pollInterval)
	reportIntervalTicker := time.NewTicker(reportInterval)
	for {
		select {
		case <-pollIntervalTicker.C:
			//runtime.ReadMemStats(&rtm)
			a.collectMetrics()
			//log.Println(time.Now().Format(time.UnixDate), "Counter update metrics: ", pollCount)
		case <-reportIntervalTicker.C:
			a.postMetrics(endpoint)
		}
	}
}
