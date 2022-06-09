package poolmemstats

import (
	"context"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/yledovskikh/devops-tpl/internal/config"
	"github.com/yledovskikh/devops-tpl/internal/serializer"
	"github.com/yledovskikh/devops-tpl/internal/storage"
)

var (
	rtm       runtime.MemStats
	poolCount int64
	metrics   []storage.Metric
)

func serializeMetrics(key string) {

	rand.Seed(time.Now().UnixNano())
	r := rand.Float64()

	metrics = []storage.Metric{
		serializer.SerializeGaugeH("Alloc", float64(rtm.Alloc), key),
		serializer.SerializeGaugeH("BuckHashSys", float64(rtm.BuckHashSys), key),
		serializer.SerializeGaugeH("Frees", float64(rtm.Frees), key),
		serializer.SerializeGaugeH("GCCPUFraction", float64(rtm.GCCPUFraction), key),
		serializer.SerializeGaugeH("GCSys", float64(rtm.GCSys), key),
		serializer.SerializeGaugeH("HeapAlloc", float64(rtm.HeapAlloc), key),
		serializer.SerializeGaugeH("HeapIdle", float64(rtm.HeapIdle), key),
		serializer.SerializeGaugeH("HeapInuse", float64(rtm.HeapInuse), key),
		serializer.SerializeGaugeH("HeapObjects", float64(rtm.HeapObjects), key),
		serializer.SerializeGaugeH("HeapReleased", float64(rtm.HeapReleased), key),
		serializer.SerializeGaugeH("HeapSys", float64(rtm.HeapSys), key),
		serializer.SerializeGaugeH("LastGC", float64(rtm.LastGC), key),
		serializer.SerializeGaugeH("Lookups", float64(rtm.Lookups), key),
		serializer.SerializeGaugeH("MCacheInuse", float64(rtm.MCacheInuse), key),
		serializer.SerializeGaugeH("MCacheSys", float64(rtm.MCacheSys), key),
		serializer.SerializeGaugeH("MSpanSys", float64(rtm.MSpanSys), key),
		serializer.SerializeGaugeH("Mallocs", float64(rtm.Mallocs), key),
		serializer.SerializeGaugeH("NextGC", float64(rtm.NextGC), key),
		serializer.SerializeGaugeH("NumForcedGC", float64(rtm.NumForcedGC), key),
		serializer.SerializeGaugeH("NextGC", float64(rtm.NumGC), key),
		serializer.SerializeGaugeH("OtherSys", float64(rtm.OtherSys), key),
		serializer.SerializeGaugeH("PauseTotalNs", float64(rtm.PauseTotalNs), key),
		serializer.SerializeGaugeH("StackInuse", float64(rtm.StackInuse), key),
		serializer.SerializeGaugeH("StackSys", float64(rtm.StackSys), key),
		serializer.SerializeGaugeH("Sys", float64(rtm.Sys), key),
		serializer.SerializeGaugeH("TotalAlloc", float64(rtm.TotalAlloc), key),
		serializer.SerializeGaugeH("MSpanInuse", float64(rtm.MSpanInuse), key),
		serializer.SerializeGaugeH("NumGC", float64(rtm.NumGC), key),
		serializer.SerializeGaugeH("RandomValue", r, key),
		serializer.SerializeCounter("PollCount", poolCount, key),
	}
}

func Exec(ctx context.Context, wg *sync.WaitGroup, agentConfig config.AgentConfig, ch chan<- *[]storage.Metric) {
	pollIntervalTicker := time.NewTicker(agentConfig.PollInterval)
	reportIntervalTicker := time.NewTicker(agentConfig.ReportInterval)

	for {
		select {
		case <-pollIntervalTicker.C:
			runtime.ReadMemStats(&rtm)
			poolCount++
			log.Info().Msg("runtime.MemStats metrics was polled")
		case <-reportIntervalTicker.C:
			serializeMetrics(agentConfig.Key)
			ch <- &metrics
			poolCount = 0
			log.Info().Msg("runtime.MemStats metrics was reported")
		case <-ctx.Done():
			log.Info().Msg("polling of runtime.MemStats metrics was stopped")
			wg.Done()
			return
		}
	}
}
