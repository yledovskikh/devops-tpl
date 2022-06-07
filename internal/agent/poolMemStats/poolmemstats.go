package poolmemstats

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/yledovskikh/devops-tpl/internal/config"
	"github.com/yledovskikh/devops-tpl/internal/hash"
	"github.com/yledovskikh/devops-tpl/internal/serializer"
	"github.com/yledovskikh/devops-tpl/internal/storage"
)

type MetricMemStats struct {
	storage storage.Storage
}

func New(storage storage.Storage) *MetricMemStats {
	return &MetricMemStats{
		storage: storage,
	}
}

func (a *MetricMemStats) pollMetrics() {

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
	log.Info().Msgf("memStats metrics was polled")
}

func (a *MetricMemStats) postBatchMetrics(key string, ch chan<- *[]storage.Metric) {

	var metrics []storage.Metric

	for mName, mValue := range a.storage.GetAllGauges() {
		var h string
		if key != "" {
			data := fmt.Sprintf("%s:gauge:%f", mName, mValue)
			h = hash.SignData(key, data)
		}

		m := serializer.SerializeGauge(mName, mValue, h)
		metrics = append(metrics, m)
	}

	for mName, mValue := range a.storage.GetAllCounters() {
		var h string
		if key != "" {
			data := fmt.Sprintf("%s:counter:%d", mName, mValue)
			h = hash.SignData(key, data)
		}
		m := serializer.SerializeCounter(mName, mValue, h)
		metrics = append(metrics, m)
	}

	ch <- &metrics

	//payloadBuf := new(bytes.Buffer)
	//if err := json.NewEncoder(payloadBuf).Encode(metrics); err != nil {
	//	log.Error().Err(err).Msg("")
	//}

}

func (a *MetricMemStats) Exec(ctx context.Context, agentConfig config.AgentConfig, ch chan<- *[]storage.Metric) {
	pollIntervalTicker := time.NewTicker(agentConfig.PollInterval)
	reportIntervalTicker := time.NewTicker(agentConfig.ReportInterval)

	for {
		select {
		case <-pollIntervalTicker.C:
			a.pollMetrics()
		case <-reportIntervalTicker.C:
			a.postBatchMetrics(agentConfig.Key, ch)
		case <-ctx.Done():
			log.Info().Msg("polling of runtime.MemStats metrics stopped")
			return
		}
	}
}
