package poolGoPsUtil

import (
	"context"
	"fmt"
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

func setMetricsGauge(metricName, key string, metricValue float64, metrics []storage.Metric) *[]storage.Metric {
	m := serializer.SerializeGaugeH(metricName, metricValue, key)
	metrics = append(metrics, m)
	return &metrics

}

func postBatchMetrics(key string, ch chan<- *[]storage.Metric) {
	var metrics []storage.Metric
	var m storage.Metric
	v, _ := mem.VirtualMemory()
	proc, _ := cpu.Percent(0, true)

	//m:=serializer.SerializeGauge("TotalMemory", float64(v.Total),"asdf")
	//a.storage.SetGauge("FreeMemory", float64(v.Free))
	//a.storage.SetGauge("CPUutilization1", c[0])
	//metrics = setMetricsGauge("TotalMemory", key, float64(v.Total), metrics)
	m = serializer.SerializeGaugeH("TotalMemory", float64(v.Total), key)
	metrics = append(metrics, m)
	m = serializer.SerializeGaugeH("FreeMemory", float64(v.Free), key)
	metrics = append(metrics, m)

	for i, p := range proc {
		m = serializer.SerializeGaugeH(fmt.Sprintf("CPUutilization%d", i), p, key)
		metrics = append(metrics, m)
	}

	ch <- &metrics
}

func Exec(ctx context.Context, agentConfig config.AgentConfig, ch chan<- *[]storage.Metric) {
	reportIntervalTicker := time.NewTicker(agentConfig.ReportInterval)

	for {
		select {
		case <-reportIntervalTicker.C:
			postBatchMetrics(agentConfig.Key, ch)
		case <-ctx.Done():
			log.Info().Msg("polling of runtime.MemStats metrics stopped")
			return
		}
	}
}
