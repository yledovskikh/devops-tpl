package storage

//type RunTimeMetrics struct {
var Counter = make(map[string]int64)
var Gauge = make(map[string]float64)

//}

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
