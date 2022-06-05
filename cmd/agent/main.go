package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yledovskikh/devops-tpl/internal/agent"
	"github.com/yledovskikh/devops-tpl/internal/config"
	"github.com/yledovskikh/devops-tpl/internal/inmemory"
)

func main() {
	s := inmemory.NewMetricStore()
	h := agent.New(s)
	agentConfig := config.GetAgentConfig()
	log.Logger = log.With().Caller().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Info().Msgf("endpoint: %s, pollInterval: %s , reportInterval: %s", agentConfig.EndPoint, agentConfig.PollInterval, agentConfig.ReportInterval)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go h.Exec(agentConfig)
	exitCode := <-signalChannel
	fmt.Println(exitCode)
	log.Info().Msgf("exit signal %s", exitCode)
}
