package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type RunTimeMetrics struct {
	counter map[string]int64
	gauge   map[string]float64
}

var m RunTimeMetrics

func (rtm *RunTimeMetrics) UpdateRTMetric(mtype string, mname string, mvalue string) error {
	switch strings.ToLower(mtype) {
	case "gauge":
		vg, err := strconv.ParseFloat(mvalue, 64)
		if err != nil {
			return errors.New("incorrect metric value")
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
			return errors.New("incorrect metric value")
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

func splitPath(sPath string) []string {
	s := strings.Split(sPath, "/")
	return s
}

func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only Post requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	//reg, err := regexp.MatchString(`^\/update\/(counter|gauge)\/\w+\/(\d+(?:\.\d+))$`, r.URL.Path)
	curl, err := regexp.MatchString(`^\/update\/(counter|gauge)\/\w+\/(\d+(?:\.\d+)?)$`, r.URL.Path)
	fmt.Println("", curl, err)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !curl {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	//if r.Header.Get("Content-Type") != "text/plain" {
	//	//http.Error(w, r.Header.Get("Content-Type"), http.StatusUnsupportedMediaType)
	//	http.Error(w, "Content-Type text/plain is required!", http.StatusUnsupportedMediaType)
	//	return
	//}

	s := splitPath(r.URL.Path)
	err = m.UpdateRTMetric(s[2], s[3], s[4])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(m)
	//fmt.Fprintln(w, m.gauge)
	//w.Write(["aa"])
	//fmt.Fprintln(w, )
	//fmt.Fprintln(w, r.URL.Path)
}
