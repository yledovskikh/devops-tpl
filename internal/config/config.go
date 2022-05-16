package config

import (
	"github.com/caarlos0/env/v6"
	"log"
	"time"
)

const pollIntervalDefault = 2 * time.Second
const reportIntervalDefault = 10 * time.Second
const serverAddresstDefault = "127.0.0.1:8080"
const schemaDefault = "http://"

type Config struct {
	ServerAddress  string        `env:"ADDRESS"`
	PollInterval   time.Duration `env:"REPORT_INTERVAL"`
	ReportInterval time.Duration `env:"POLL_INTERVAL"`
}

func AgentConfig() (string, time.Duration, time.Duration) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	validate(&cfg)
	return schemaDefault + cfg.ServerAddress, cfg.PollInterval, cfg.ReportInterval
}

func ServerConfig() string {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	validate(&cfg)
	return cfg.ServerAddress
}

func validate(cfg *Config) {
	if cfg.ServerAddress == "" {
		cfg.ServerAddress = serverAddresstDefault
	}
	if cfg.PollInterval == 0 {
		cfg.PollInterval = pollIntervalDefault
	}
	if cfg.ReportInterval == 0 {
		cfg.ReportInterval = reportIntervalDefault
	}
}
