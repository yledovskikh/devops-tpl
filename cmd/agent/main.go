package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yledovskikh/devops-tpl/internal/agent/poolGoPsUtil"
	"github.com/yledovskikh/devops-tpl/internal/agent/poolMemStats"
	"github.com/yledovskikh/devops-tpl/internal/agent/reportMetrics"
	"github.com/yledovskikh/devops-tpl/internal/inmemory"
	"github.com/yledovskikh/devops-tpl/internal/storage"

	"github.com/yledovskikh/devops-tpl/internal/config"
)

func main() {
	s := inmemory.NewMetricStore()
	h := poolmemstats.New(s)
	ch := make(chan *[]storage.Metric)

	log.Logger = log.With().Caller().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	ctx, cancel := context.WithCancel(context.Background())

	agentConfig := config.GetAgentConfig()
	log.Info().Msgf("endpoint: %s, pollInterval: %s , reportInterval: %s", agentConfig.EndPoint, agentConfig.PollInterval, agentConfig.ReportInterval)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go h.Exec(ctx, agentConfig, ch)
	go poolgopsutil.Exec(ctx, agentConfig, ch)
	go reportmetrics.Exec(ctx, agentConfig.EndPoint, ch)

	exitCode := <-signalChannel
	cancel()
	fmt.Println(exitCode)
	log.Info().Msgf("exit signal %s", exitCode)
}
