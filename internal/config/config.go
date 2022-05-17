package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	serverAddressDefault  = "127.0.0.1:8080"
	pollIntervalDefault   = 2 * time.Second
	reportIntervalDefault = 10 * time.Second

	storeIntervalDefault = 300 * time.Second
	storeFileDefault     = "/tmp/devops-metrics-db.json"
	restoreDefault       = true
)

type AgentConfig struct {
	ServerAddress  string `env:"ADDRESS"`
	EndPoint       string
	PollInterval   time.Duration `env:"REPORT_INTERVAL"`
	ReportInterval time.Duration `env:"POLL_INTERVAL"`
}

type ServerConfig struct {
	ServerAddress string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool
}

var pollInterval time.Duration
var reportInterval time.Duration
var serverAddress string

func validateAgentConfig(cfg *AgentConfig) {
	if cfg.ServerAddress == "" {
		cfg.ServerAddress = serverAddress
	}
	if cfg.PollInterval == 0 {
		cfg.PollInterval = pollInterval
	}
	if cfg.ReportInterval == 0 {
		cfg.ReportInterval = reportInterval
	}

	cfg.EndPoint = "https://" + cfg.ServerAddress

}

func validateServerConfig(cfg *ServerConfig) {
	restoreEnv := os.Getenv("RESTORE")
	if cfg.ServerAddress == "" {
		cfg.ServerAddress = serverAddressDefault
	}
	if cfg.StoreFile == "" {
		cfg.StoreFile = storeFileDefault
	}
	if cfg.StoreInterval == 0 {
		cfg.StoreInterval = storeIntervalDefault
	}

	if restoreEnv == "" {
		cfg.Restore = restoreDefault
	} else {
		var err error
		cfg.Restore, err = strconv.ParseBool(restoreEnv)
		if err != nil {
			cfg.Restore = restoreDefault
		}
	}

}

func GetAgentConfig() AgentConfig {
	var cfg AgentConfig

	flag.StringVar(&serverAddress, "a", serverAddressDefault, "server address")
	flag.DurationVar(&pollInterval, "p", pollIntervalDefault, "poll metrics interval")
	flag.DurationVar(&reportInterval, "r", reportIntervalDefault, "report metric interval")

	flag.Parse()
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	validateAgentConfig(&cfg)
	return cfg
}

func GetServerConfig() ServerConfig {
	var cfg ServerConfig
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	validateServerConfig(&cfg)
	return cfg
}
