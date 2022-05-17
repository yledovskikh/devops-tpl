package config

import (
	"github.com/caarlos0/env/v6"
	"log"
	"strconv"
	"time"
)

const (
	serverAddressDefault   = "127.0.0.1:8080"
	pollIntervalDefault    = 2 * time.Second
	reportIntervalDefault  = 10 * time.Second
	serverSSLEnableDefault = false

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
	RestoreEnv    string        `env:"RESTORE"`
	Restore       bool
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

func validateServerConfig(cfg *ServerConfig) {
	if cfg.ServerAddress == "" {
		cfg.ServerAddress = serverAddressDefault
	}
	if cfg.StoreFile == "" {
		cfg.StoreFile = storeFileDefault
	}
	if cfg.StoreInterval == 0 {
		cfg.StoreInterval = storeIntervalDefault
	}

	if cfg.RestoreEnv == "" {
		cfg.Restore = restoreDefault
	} else {
		var err error
		cfg.Restore, err = strconv.ParseBool(cfg.RestoreEnv)
		if err != nil {
			cfg.Restore = restoreDefault
		}
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
	validateServerConfig(&cfg)
	return cfg
}
