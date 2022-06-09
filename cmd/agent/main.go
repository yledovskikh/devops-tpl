package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yledovskikh/devops-tpl/internal/agent/poolGoPsUtil"
	"github.com/yledovskikh/devops-tpl/internal/agent/poolMemStats"
	"github.com/yledovskikh/devops-tpl/internal/agent/reportMetrics"
	"github.com/yledovskikh/devops-tpl/internal/storage"

	"github.com/yledovskikh/devops-tpl/internal/config"
)

func main() {
	ch := make(chan *[]storage.Metric)

	log.Logger = log.With().Caller().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	agentConfig := config.GetAgentConfig()
	log.Info().Msgf("endpoint: %s, pollInterval: %s , reportInterval: %s", agentConfig.EndPoint, agentConfig.PollInterval, agentConfig.ReportInterval)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	wg.Add(3)
	go poolmemstats.Exec(ctx, &wg, agentConfig, ch)
	go poolgopsutil.Exec(ctx, &wg, agentConfig, ch)
	go reportmetrics.Exec(ctx, &wg, agentConfig.EndPoint, ch)

	exitCode := <-signalChannel
	cancel()
	wg.Wait()
	fmt.Println(exitCode)
	log.Info().Msgf("exit signal %s", exitCode)
}
