package config

import (
	"github.com/caarlos0/env/v6"
	"log"
	"time"
)

type AgentConfig struct {
	ServerAddress  string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	EndPoint       string
	PollInterval   time.Duration `env:"REPORT_INTERVAL" envDefault:"2s"`
	ReportInterval time.Duration `env:"POLL_INTERVAL" envDefault:"10s"`
}

type ServerConfig struct {
	ServerAddress string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
}

func GetAgentConfig() AgentConfig {
	var cfg AgentConfig
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	//validateAgent(&cfg)
	cfg.EndPoint = "http://" + cfg.ServerAddress
	//return cfg.EndPoint, cfg.PollInterval, cfg.ReportInterval
	return cfg
}

func GetServerConfig() ServerConfig {
	var cfg ServerConfig
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	//validateServer(&cfg)
	return cfg
}

//func validateAgent(cfg *AgentConfig) {
//	if cfg.ServerAddress == "" {
//		cfg.ServerAddress = serverAddresstDefault
//	}
//	if cfg.PollInterval == 0 {
//		cfg.PollInterval = pollIntervalDefault
//	}
//	if cfg.ReportInterval == 0 {
//		cfg.ReportInterval = reportIntervalDefault
//	}
//}
//
//func validateServer(cfg *ServerConfig) {
//	if cfg.ServerAddress == "" {
//		cfg.ServerAddress = serverAddresstDefault
//	}
//}
