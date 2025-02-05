package main

import (
	"time"

	"github.com/hashicorp/go-metrics"
	"github.com/hashicorp/go-metrics/prometheus"
)

func main() {
	host := "localhost:9092"
	sink, err := prometheus.NewPrometheusPushSink(host, time.Second, "pushtest")
	if err != nil {
		panic(err)
	}
	metricsConf := metrics.DefaultConfig("default")
	metricsConf.EnableHostnameLabel = true
	metrics.NewGlobal(metricsConf, sink)

	metrics.SetGauge([]string{"one", "two"}, 42)

	for {
		time.Sleep(10 * time.Second)
		metrics.IncrCounter([]string{"three", "four"}, 1)
	}

}
