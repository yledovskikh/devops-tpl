package agent

import (
	"bytes"
	"fmt"
	"github.com/yledovskikh/devops-tpl/internal/serializer"
	"github.com/yledovskikh/devops-tpl/internal/storage"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

//type metric struct {
//	name       string
//	metricType string
//	value      string
//}

type Agent struct {
	storage storage.Storage
}

func New(storage storage.Storage) *Agent {
	return &Agent{
		storage: storage,
	}
}

func (a *Agent) Put(metricType string, metricName string, metricValue string) {
	err := a.storage.Put(metricType, metricName, metricValue)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (a *Agent) collectMetrics(rtm runtime.MemStats, pollCount int64, randomValue float64) {

	a.Put("gauge", "Alloc", strconv.FormatUint(rtm.Alloc, 10))
	a.Put("gauge", "BuckHashSys", strconv.FormatUint(rtm.BuckHashSys, 10))
	a.Put("gauge", "Frees", strconv.FormatUint(rtm.Frees, 10))
	a.Put("gauge", "GCCPUFraction", strconv.FormatFloat(rtm.GCCPUFraction, 'f', -1, 64))
	a.Put("gauge", "HeapAlloc", strconv.FormatUint(rtm.HeapAlloc, 10))
	a.Put("gauge", "HeapInuse", strconv.FormatUint(rtm.HeapInuse, 10))
	a.Put("gauge", "Alloc", strconv.FormatUint(rtm.Alloc, 10))
	//...
	////Custom metrics
	a.Put("counter", "PollCount", strconv.FormatInt(pollCount, 10))
	a.Put("gauge", "RandomValue", strconv.FormatFloat(randomValue, 'f', -1, 64))

	//
	//{ "gauge", "HeapObjects", strconv.FormatUint(rtm.HeapObjects, 10)},
	//{"gauge","HeapReleased",  strconv.FormatUint(rtm.HeapReleased, 10)},
	//{ "gauge","HeapSys", strconv.FormatUint(rtm.HeapSys, 10)},
	//{"gauge", "LastGC",  strconv.FormatUint(rtm.LastGC, 10)},
	//{ "gauge","Lookups", strconv.FormatUint(rtm.Lookups, 10)},
	//{"gauge", "MCacheInuse",  strconv.FormatUint(rtm.MCacheInuse, 10)},
	//{ "gauge", "MCacheSys", strconv.FormatUint(rtm.MCacheSys, 10)},
	//{"gauge", "MSpanInuse",  strconv.FormatUint(rtm.MSpanInuse, 10)},
	//{ "gauge","MSpanSys", strconv.FormatUint(rtm.MSpanSys, 10)},
	//{"gauge", "Mallocs",  strconv.FormatUint(rtm.Mallocs, 10)},
	//{ "gauge","NextGC", strconv.FormatUint(rtm.NextGC, 10)},
	//{"gauge","NumForcedGC",  strconv.FormatUint(uint64(rtm.NumForcedGC), 10)},
	//{ "gauge","NumGC", strconv.FormatUint(uint64(rtm.NumGC), 10)},
	//{"gauge", "OtherSys",  strconv.FormatUint(rtm.OtherSys, 10)},
	//{ "gauge", "PauseTotalNs", strconv.FormatUint(rtm.PauseTotalNs, 10)},
	//{"gauge", "StackInuse",  strconv.FormatUint(rtm.StackInuse, 10)},
	//{ "gauge", "StackSys", strconv.FormatUint(rtm.StackSys, 10)},
	//{"gauge", "Sys",  strconv.FormatUint(rtm.Sys, 10)},
	//{ "gauge", "TotalAlloc", strconv.FormatUint(rtm.TotalAlloc, 10)},

	//m := []metric{
	//	{"Alloc", "gauge", strconv.FormatUint(rtm.Alloc, 10)},
	//	{"BuckHashSys", "gauge", strconv.FormatUint(rtm.BuckHashSys, 10)},
	//	{"Frees", "gauge", strconv.FormatUint(rtm.Frees, 10)},
	//	{"GCCPUFraction", "gauge", strconv.FormatFloat(rtm.GCCPUFraction, 'f', -1, 64)},
	//	{"GCSys", "gauge", strconv.FormatUint(rtm.GCSys, 10)},
	//	{"HeapAlloc", "gauge", strconv.FormatUint(rtm.HeapAlloc, 10)},
	//	{"HeapIdle", "gauge", strconv.FormatUint(rtm.HeapIdle, 10)},
	//	{"HeapInuse", "gauge", strconv.FormatUint(rtm.HeapInuse, 10)},
	//	{"HeapObjects", "gauge", strconv.FormatUint(rtm.HeapObjects, 10)},
	//	{"HeapReleased", "gauge", strconv.FormatUint(rtm.HeapReleased, 10)},
	//	{"HeapSys", "gauge", strconv.FormatUint(rtm.HeapSys, 10)},
	//	{"LastGC", "gauge", strconv.FormatUint(rtm.LastGC, 10)},
	//	{"Lookups", "gauge", strconv.FormatUint(rtm.Lookups, 10)},
	//	{"MCacheInuse", "gauge", strconv.FormatUint(rtm.MCacheInuse, 10)},
	//	{"MCacheSys", "gauge", strconv.FormatUint(rtm.MCacheSys, 10)},
	//	{"MSpanInuse", "gauge", strconv.FormatUint(rtm.MSpanInuse, 10)},
	//	{"MSpanSys", "gauge", strconv.FormatUint(rtm.MSpanSys, 10)},
	//	{"Mallocs", "gauge", strconv.FormatUint(rtm.Mallocs, 10)},
	//	{"NextGC", "gauge", strconv.FormatUint(rtm.NextGC, 10)},
	//	{"NumForcedGC", "gauge", strconv.FormatUint(uint64(rtm.NumForcedGC), 10)},
	//	{"NumGC", "gauge", strconv.FormatUint(uint64(rtm.NumGC), 10)},
	//	{"OtherSys", "gauge", strconv.FormatUint(rtm.OtherSys, 10)},
	//	{"PauseTotalNs", "gauge", strconv.FormatUint(rtm.PauseTotalNs, 10)},
	//	{"StackInuse", "gauge", strconv.FormatUint(rtm.StackInuse, 10)},
	//	{"StackSys", "gauge", strconv.FormatUint(rtm.StackSys, 10)},
	//	{"Sys", "gauge", strconv.FormatUint(rtm.Sys, 10)},
	//	{"TotalAlloc", "gauge", strconv.FormatUint(rtm.TotalAlloc, 10)},
	//	//Custom metrics
	//	{"PollCount", "counter", strconv.FormatInt(pollCount, 10)},
	//	{"RandomValue", "gauge", strconv.FormatFloat(randomValue, 'f', -1, 64)},
	//}
	//return m
}
func send2server(url string, body []byte) error {

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
	fmt.Println(url)
	return nil
}

//TODO join func postGauges and postCounter

func (a *Agent) postGauges(updateMetricURL string) {
	fmt.Println(time.Now().Format(time.UnixDate), "Push Gauges metrics:")
	for mName, mValue := range a.storage.GetAllGauges() {
		url := updateMetricURL
		//url := updateMetricURL + "/gauge/" + mName + "/" + fmt.Sprintf("%f", mValue)

		body, err := serializer.EncodingMetricGauge(mName, mValue)
		//TODO remove debug info
		fmt.Printf(string(body))

		if err != nil {
			fmt.Println(err)
		}

		if err = send2server(url, body); err != nil {
			fmt.Println(err)
		}
	}
}

func (a *Agent) postCounter(updateMetricURL string) {
	fmt.Println(time.Now().Format(time.UnixDate), "Push Counter metrics:")
	for mName, mValue := range a.storage.GetAllCounters() {
		url := updateMetricURL
		body, err := serializer.EncodingMetricCounter(mName, mValue)
		//TODO remove debug info
		fmt.Printf(string(body))

		if err != nil {
			fmt.Println(err)
		}

		if err = send2server(url, body); err != nil {
			fmt.Println(err)
		}
	}
}

func (a *Agent) Exec(pollInterval time.Duration, reportInterval time.Duration, updateMetricURL string) {
	var rtm runtime.MemStats
	var pollCount int64 = 0
	pollIntervalTicker := time.NewTicker(pollInterval)
	reportIntervalTicker := time.NewTicker(reportInterval)
	rand.Seed(time.Now().UnixNano())
	for {
		select {
		case <-pollIntervalTicker.C:
			pollCount++
			runtime.ReadMemStats(&rtm)
			fmt.Println(time.Now().Format(time.UnixDate), "Counter update metrics: ", pollCount)
		case <-reportIntervalTicker.C:
			r := rand.Float64()
			a.collectMetrics(rtm, pollCount, r)
			a.postGauges(updateMetricURL)
			a.postCounter(updateMetricURL)
			pollCount = 0
		}
	}
}
