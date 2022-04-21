package main

import (
	"errors"
	"strconv"
	"strings"
)

//type gauge float64
//type counter int64

//type RunTimeMetrics struct {
//	Alloc         gauge
//	BuckHashSys   gauge
//	Frees         gauge
//	GCCPUFraction gauge
//	GCSys         gauge
//	HeapAlloc     gauge
//	HeapIdle      gauge
//	HeapInuse     gauge
//	HeapObjects   gauge
//	HeapReleased  gauge
//	HeapSys       gauge
//	LastGC        gauge
//	Lookups       gauge
//	MCacheInuse   gauge
//	MCacheSys     gauge
//	MSpanInuse    gauge
//	MSpanSys      gauge
//	Mallocs       gauge
//	NextGC        gauge
//	NumForcedGC   gauge
//	NumGC         gauge
//	OtherSys      gauge
//	PauseTotalNs  gauge
//	StackInuse    gauge
//	StackSys      gauge
//	Sys           gauge
//	TotalAlloc    gauge
//	RandomValue   gauge
//	PollCount     counter
//}

//func (rtm *RunTimeMetrics) setMgauge(m string, v gauge) {
//	rtm.mgauge[m]=v
//}

//func string2gauge(v string) (float64,error) {
//	vg,err := strconv.ParseFloat(v,64)
//	return vg,err
//}

type RunTimeMetrics struct {
	counter map[string]int64
	gauge   map[string]float64
}

func (rtm *RunTimeMetrics) UpdateRTMetric(mtype string, mname string, mvalue string) error {
	switch strings.ToLower(mtype) {
	case "gauge":
		vg, err := strconv.ParseFloat(mvalue, 64)
		if err != nil {
			return errors.New("Incorrect metric value")
		}
		//fmt.Println(mname, vg)
		//rtm.gauge[mname] = vg
		if rtm.gauge == nil {
			rtm.gauge = make(map[string]float64)
		}
		rtm.gauge[mname] = vg
	case "counter":
		vg, err := strconv.ParseInt(mvalue, 10, 64)
		if err != nil {
			return errors.New("Incorrect metric value")
		}
		if rtm.counter == nil {
			rtm.counter = make(map[string]int64)
		}
		rtm.counter[mname] += vg
	default:
		return errors.New("incorrect type (expected gauge or counter)")
	}
	return nil
}
