package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yledovskikh/devops-tpl/internal/agent"
	"github.com/yledovskikh/devops-tpl/internal/config"
	"github.com/yledovskikh/devops-tpl/internal/storage"
)

func main() {
	s := storage.NewMetricStore()
	h := agent.New(s)
	agentConfig := config.GetAgentConfig()
	log.Printf("endpoint: %s, pollInterval: %s , reportInterval: %s", agentConfig.EndPoint, agentConfig.PollInterval, agentConfig.ReportInterval)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go h.Exec(agentConfig)
	exitCode := <-signalChannel
	fmt.Println(exitCode)

}
