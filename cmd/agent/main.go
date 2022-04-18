package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"
	//"time"
)

const (
	endpoint       = "http://localhost:8080"
	contextURL     = "update"
	pollInterval   = 2
	reportInterval = 10
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

func postMetrics(m []metric) {
	for _, value := range m {
		url := endpoint + "/" + contextURL + "/" + value.metricType + "/" + value.name + "/" + value.value
		fmt.Println(url)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		request.Header.Set("Content-Type", "text/plain; charset=UTF-8")
		if err != nil {
			fmt.Println(err.Error())
		}
		client := &http.Client{}
		client.Do(request)
		response, err := client.Do(request)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(url, response.StatusCode)
	}

}

func main() {
	var m []metric
	var rtm runtime.MemStats
	rand.Seed(time.Now().UnixNano())
	var pollCount int64 = 0
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT)
	exitChan := make(chan int)
	go func() {

		for {
			sig := <-signalChannel
			switch sig {
			case syscall.SIGTERM:
				fmt.Println("sigterm")
				exitChan <- 0
			case syscall.SIGINT:
				fmt.Println("sigint")
				exitChan <- 0
			case syscall.SIGQUIT:
				fmt.Println("sigquit")
				exitChan <- 0
			default:
				fmt.Println("Unknown signal.")
				exitChan <- 1
			}
		}

	}()

	go func() {
		for {
			pollCount++
			runtime.ReadMemStats(&rtm)
			r := rand.Float64()
			m = collectMetrics(rtm, pollCount, r)
			fmt.Println(pollCount)
			if pollCount%reportInterval == 0 {
				postMetrics(m)
			}
			time.Sleep(pollInterval * time.Second)
		}
	}()
	exitCode := <-exitChan
	os.Exit(exitCode)
}
