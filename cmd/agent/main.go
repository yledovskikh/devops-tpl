package main

import (
	"github.com/yledovskikh/devops-tpl/internal/agent"
	"os"
	"time"
)

const (
	endpoint       = "http://localhost:8080"
	contextURL     = "update"
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func main() {

	exitChan := make(chan int)
	go agent.TerminateAgent(exitChan)
	go agent.RefreshMetrics(pollInterval, reportInterval, endpoint, contextURL)
	exitCode := <-exitChan
	os.Exit(exitCode)
}
