package main

import (
	"fmt"
	"github.com/yledovskikh/devops-tpl/internal/agent"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	//endpoint       	= "http://localhost:8080"
	//contextURL     	= "update"
	updateMetricURL = "http://localhost:8080/update"
	pollInterval    = 2 * time.Second
	reportInterval  = 10 * time.Second
)

func main() {

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go agent.RefreshMetrics(pollInterval, reportInterval, updateMetricURL)
	exitCode := <-signalChannel
	fmt.Println(exitCode)

}
