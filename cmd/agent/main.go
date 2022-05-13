package main

import (
	"fmt"
	"github.com/yledovskikh/devops-tpl/internal/agent"
	"github.com/yledovskikh/devops-tpl/internal/storage"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	updateMetricURL = "http://localhost:8080/update"
	pollInterval    = 2 * time.Second
	reportInterval  = 10 * time.Second
)

func main() {
	s := storage.NewMetricStore()
	h := agent.New(s)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go h.Exec(pollInterval, reportInterval, updateMetricURL)
	exitCode := <-signalChannel
	fmt.Println(exitCode)

}
