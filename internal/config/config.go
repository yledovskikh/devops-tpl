package config

import (
	"github.com/caarlos0/env/v6"
	"log"
	"time"
)

const (
	serverAddressDefault   = "127.0.0.1:8080"
	pollIntervalDefault    = 2 * time.Second
	reportIntervalDefault  = 10 * time.Second
	serverSSLEnableDefault = false
)

type AgentConfig struct {
	ServerAddress  string `env:"ADDRESS"`
	EndPoint       string
	PollInterval   time.Duration `env:"REPORT_INTERVAL"`
	ReportInterval time.Duration `env:"POLL_INTERVAL"`
}

type ServerConfig struct {
	ServerAddress string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
}

//func init() {
//
//}

func validateAgentConfig(cfg *AgentConfig) {
	if cfg.ServerAddress == "" {
		cfg.ServerAddress = serverAddressDefault
	}
	if cfg.PollInterval == 0 {
		cfg.PollInterval = pollIntervalDefault
	}
	if cfg.ReportInterval == 0 {
		cfg.ReportInterval = reportIntervalDefault
	}
	if serverSSLEnableDefault {
		cfg.EndPoint = "https://" + cfg.ServerAddress
	} else {
		cfg.EndPoint = "http://" + cfg.ServerAddress
	}

}

func GetAgentConfig() AgentConfig {
	var cfg AgentConfig
	err := env.Parse(&cfg)
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
	return cfg
}
