package main

import (
	"fmt"
	"github.com/yledovskikh/devops-tpl/internal/agent"
	"github.com/yledovskikh/devops-tpl/internal/config"
	"github.com/yledovskikh/devops-tpl/internal/storage"
	"log"
	"os"
	"os/signal"
	"syscall"
)

//const (
//	updateMetricURL = "http://localhost:8080/update/"
//	pollInterval    = 2 * time.Second
//	reportInterval  = 10 * time.Second
//)
//
//type Config struct {
//	Files        []string      `env:"FILES" envSeparator:":"`
//	Home         string        `env:"HOME"`
//	TaskDuration time.Duration `env:"TASK_DURATION,required"`
//
//	endpoint       string        `env:"ADDRESS"`
//	pollInterval   time.Duration `env:"REPORT_INTERVAL"`
//	reportInterval time.Duration `env:"POLL_INTERVAL"`
//}

func main() {

	//var cfg Config
	//err := env.Parse(&cfg)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//log.Println(cfg)

	//cfg.endpoint = os.Getenv("ADDRESS")
	//log.Println("endpoint, pollInterval, reportInterval:", cfg.endpoint)
	//log.Println("endpoint, pollInterval, reportInterval:", cfg.endpoint, cfg.pollInterval, cfg.reportInterval)

	s := storage.NewMetricStore()
	h := agent.New(s)
	endpoint, pollInterval, reportInterval := config.AgentConfig()
	log.Printf("endpoint: %s, pollInterval: %s , reportInterval: %s", endpoint, pollInterval, reportInterval)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go h.Exec(endpoint, pollInterval, reportInterval)
	exitCode := <-signalChannel
	fmt.Println(exitCode)

}
