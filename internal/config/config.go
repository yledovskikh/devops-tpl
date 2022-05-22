package config

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/caarlos0/env/v6"
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
	EndPoint       string
	PollInterval   time.Duration
	ReportInterval time.Duration
}

type AgentConfigEnv struct {
	EndPoint       string        `env:"ADDRESS"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
}

type ServerConfig struct {
	ServerAddress string
	StoreInterval time.Duration
	StoreFile     string
	Restore       bool
}

type ServerConfigEnv struct {
	ServerAddress string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool
}

var pollInterval time.Duration
var reportInterval time.Duration
var storeInterval time.Duration
var serverAddress string
var storeFile string
var restoreServer bool

func validateAgentConfig(cfg *AgentConfig, cEnv *AgentConfigEnv) {

	cfg.EndPoint = cEnv.EndPoint
	if cEnv.EndPoint == "" {
		cfg.EndPoint = serverAddress
	}

	cfg.PollInterval = cEnv.PollInterval
	if cEnv.PollInterval == time.Duration(0) {
		cfg.PollInterval = pollInterval
	}

	cfg.ReportInterval = cEnv.ReportInterval
	if cEnv.ReportInterval == time.Duration(0) {
		cfg.ReportInterval = reportInterval
	}

	cfg.EndPoint = "http://" + cfg.EndPoint
}

func validateServerConfig(cfg *ServerConfig, cEnv *ServerConfigEnv) {

	cfg.ServerAddress = cEnv.ServerAddress
	if cEnv.ServerAddress == "" {
		cfg.ServerAddress = serverAddress
	}
	cfg.StoreFile = cEnv.StoreFile
	if cEnv.StoreFile == "" {
		cfg.StoreFile = storeFile
	}
	cfg.StoreInterval = cEnv.StoreInterval
	if cfg.StoreInterval == time.Duration(0) {
		cfg.StoreInterval = storeInterval
	}

	restoreEnv := os.Getenv("RESTORE")
	if restoreEnv == "" {
		cfg.Restore = restoreServer
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
	var cEnv AgentConfigEnv

	flag.StringVar(&serverAddress, "a", serverAddressDefault, "server address")
	flag.DurationVar(&pollInterval, "p", pollIntervalDefault, "poll metrics interval")
	flag.DurationVar(&reportInterval, "r", reportIntervalDefault, "report metric interval")

	flag.Parse()
	err := env.Parse(&cEnv)
	if err != nil {
		log.Fatal(err)
	}

	validateAgentConfig(&cfg, &cEnv)
	return cfg
}

func GetServerConfig() ServerConfig {
	var cfg ServerConfig
	var cEnv ServerConfigEnv

	flag.StringVar(&serverAddress, "a", serverAddressDefault, "server address")
	flag.DurationVar(&storeInterval, "i", storeIntervalDefault, "dump metrics to file interval")
	flag.StringVar(&storeFile, "f", storeFileDefault, "dump file name")
	flag.BoolVar(&restoreServer, "r", restoreDefault, "restore metrics from file")

	flag.Parse()
	err := env.Parse(&cEnv)
	if err != nil {
		log.Println(err)
	}

	validateServerConfig(&cfg, &cEnv)
	return cfg
}
