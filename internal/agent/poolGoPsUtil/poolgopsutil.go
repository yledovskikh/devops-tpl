package poolgopsutil

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/yledovskikh/devops-tpl/internal/config"
	"github.com/yledovskikh/devops-tpl/internal/serializer"
	"github.com/yledovskikh/devops-tpl/internal/storage"
)

//type MetricGoPsUtl struct {
//	storage storage.Storage
//}
//
//func New(storage storage.Storage) *MetricGoPsUtl {
//	return &MetricGoPsUtl{
//		storage: storage,
//	}
//}

func postBatchMetrics(key string, ch chan<- *[]storage.Metric) {
	//var metrics []storage.Metric
	var m storage.Metric
	v, _ := mem.VirtualMemory()
	proc, _ := cpu.Percent(0, true)

	metrics := []storage.Metric{
		serializer.SerializeGaugeH("TotalMemory", float64(v.Total), key),
		serializer.SerializeGaugeH("FreeMemory", float64(v.Free), key),
	}

	for i, p := range proc {
		m = serializer.SerializeGaugeH(fmt.Sprintf("CPUutilization%d", i), p, key)
		metrics = append(metrics, m)
	}

	ch <- &metrics
}

func Exec(ctx context.Context, wg *sync.WaitGroup, agentConfig config.AgentConfig, ch chan<- *[]storage.Metric) {
	reportIntervalTicker := time.NewTicker(agentConfig.ReportInterval)

	for {
		select {
		case <-reportIntervalTicker.C:
			postBatchMetrics(agentConfig.Key, ch)
			log.Info().Msg("GoPsUtl metrics was polled")
		case <-ctx.Done():
			log.Info().Msg("polling of GoPsUtl metrics was stopped")
			wg.Done()
			return
		}
	}
}
