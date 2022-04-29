package agent

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

type metric struct {
	name       string
	metricType string
	value      string
}

func collectMetrics(rtm runtime.MemStats, pollCount int64, randomValue float64) []metric {
	m := []metric{
		{"Alloc", "gauge", strconv.FormatUint(rtm.Alloc, 10)},
		{"BuckHashSys", "gauge", strconv.FormatUint(rtm.BuckHashSys, 10)},
		{"Frees", "gauge", strconv.FormatUint(rtm.Frees, 10)},
		{"GCCPUFraction", "gauge", strconv.FormatFloat(rtm.GCCPUFraction, 'f', -1, 64)},
		{"GCSys", "gauge", strconv.FormatUint(rtm.GCSys, 10)},
		{"HeapAlloc", "gauge", strconv.FormatUint(rtm.HeapAlloc, 10)},
		{"HeapIdle", "gauge", strconv.FormatUint(rtm.HeapIdle, 10)},
		{"HeapInuse", "gauge", strconv.FormatUint(rtm.HeapInuse, 10)},
		{"HeapObjects", "gauge", strconv.FormatUint(rtm.HeapObjects, 10)},
		{"HeapReleased", "gauge", strconv.FormatUint(rtm.HeapReleased, 10)},
		{"HeapSys", "gauge", strconv.FormatUint(rtm.HeapSys, 10)},
		{"LastGC", "gauge", strconv.FormatUint(rtm.LastGC, 10)},
		{"Lookups", "gauge", strconv.FormatUint(rtm.Lookups, 10)},
		{"MCacheInuse", "gauge", strconv.FormatUint(rtm.MCacheInuse, 10)},
		{"MCacheSys", "gauge", strconv.FormatUint(rtm.MCacheSys, 10)},
		{"MSpanInuse", "gauge", strconv.FormatUint(rtm.MSpanInuse, 10)},
		{"MSpanSys", "gauge", strconv.FormatUint(rtm.MSpanSys, 10)},
		{"Mallocs", "gauge", strconv.FormatUint(rtm.Mallocs, 10)},
		{"NextGC", "gauge", strconv.FormatUint(rtm.NextGC, 10)},
		{"NumForcedGC", "gauge", strconv.FormatUint(uint64(rtm.NumForcedGC), 10)},
		{"NumGC", "gauge", strconv.FormatUint(uint64(rtm.NumGC), 10)},
		{"OtherSys", "gauge", strconv.FormatUint(rtm.OtherSys, 10)},
		{"PauseTotalNs", "gauge", strconv.FormatUint(rtm.PauseTotalNs, 10)},
		{"StackInuse", "gauge", strconv.FormatUint(rtm.StackInuse, 10)},
		{"StackSys", "gauge", strconv.FormatUint(rtm.StackSys, 10)},
		{"Sys", "gauge", strconv.FormatUint(rtm.Sys, 10)},
		{"TotalAlloc", "gauge", strconv.FormatUint(rtm.TotalAlloc, 10)},
		//Custom metrics
		{"PollCount", "counter", strconv.FormatInt(pollCount, 10)},
		{"RandomValue", "gauge", strconv.FormatFloat(randomValue, 'f', -1, 64)},
	}
	return m
}

func postMetrics(m []metric, endpoint string, contextURL string) {
	fmt.Println(time.Now().Format(time.UnixDate), "Push metrics:")
	for _, value := range m {
		url := endpoint + "/" + contextURL + "/" + value.metricType + "/" + value.name + "/" + value.value

		response, err := http.Post(url, "text/plain", nil)
		if err != nil {
			fmt.Println(err.Error())
		}
		err = response.Body.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func RefreshMetrics(pollInterval time.Duration, reportInterval time.Duration, endpoint string, contextURL string) {
	var m []metric
	var rtm runtime.MemStats
	var pollCount int64 = 0
	pollIntervalTicker := time.NewTicker(pollInterval)
	reportIntervalTicker := time.NewTicker(reportInterval)
	rand.Seed(time.Now().UnixNano())
	for {
		select {
		case <-pollIntervalTicker.C:
			r := rand.Float64()
			pollCount++
			runtime.ReadMemStats(&rtm)
			m = collectMetrics(rtm, pollCount, r)
			fmt.Println(time.Now().Format(time.UnixDate), "Counter update metrics: ", pollCount)
		case <-reportIntervalTicker.C:
			postMetrics(m, endpoint, contextURL)
		}
	}
}

//func TerminateAgent(exitChan chan int) {
//	signalChannel := make(chan os.Signal, 1)
//	//TODO refactor: think about how to optimize this structure
//	signal.Notify(signalChannel,
//		syscall.SIGTERM,
//		syscall.SIGINT,
//		syscall.SIGQUIT)
//	for {
//		sig := <-signalChannel
//		switch sig {
//		case syscall.SIGTERM:
//			fmt.Println("sigterm")
//			exitChan <- 0
//		case syscall.SIGINT:
//			fmt.Println("sigint")
//			exitChan <- 0
//		case syscall.SIGQUIT:
//			fmt.Println("sigquit")
//			exitChan <- 0
//		default:
//			fmt.Println("Unknown signal.")
//			exitChan <- 1
//		}
//	}
//}
