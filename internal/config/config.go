package config

import (
	"flag"
	"log"
	"os"
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
	Key            string
}

type AgentConfigEnv struct {
	EndPoint       string        `env:"ADDRESS"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	Key            string        `env:"KEY"`
}

type ServerConfig struct {
	ServerAddress string
	StoreInterval time.Duration
	StoreFile     string
	Restore       bool
	Key           string
	DatabaseDSN   string
}

type ServerConfigEnv struct {
	ServerAddress string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	Key           string        `env:"KEY"`
	DatabaseDSN   string        `env:"DATABASE_DSN"`
}

func validateAgentConfig(cfg *AgentConfig, cEnv *AgentConfigEnv) {

	if cEnv.EndPoint != "" {
		cfg.EndPoint = cEnv.EndPoint
	}

	if cEnv.Key != "" {
		cfg.Key = cEnv.Key
	}

	//Переделал проверку условия cEnv.ReportInterval != time.Duration(0)
	//т.к. можно выставить переменную ОС - export REPORT_INTERVAL=0s
	//по этой же причине не проверяю тип time.Duration и для других переменных

	if os.Getenv("REPORT_INTERVAL") != "" {
		cfg.ReportInterval = cEnv.ReportInterval
	}

	if os.Getenv("POLL_INTERVAL") != "" {
		cfg.PollInterval = cEnv.PollInterval
	}

	cfg.EndPoint = "http://" + cfg.EndPoint
}

func validateServerConfig(cfg *ServerConfig, cEnv *ServerConfigEnv) {

	if cEnv.ServerAddress != "" {
		cfg.ServerAddress = cEnv.ServerAddress
	}

	if cEnv.Key != "" {
		cfg.Key = cEnv.Key
	}

	if cEnv.StoreFile != "" {
		cfg.StoreFile = cEnv.StoreFile
	}

	if os.Getenv("STORE_INTERVAL") != "" {
		cfg.StoreInterval = cEnv.StoreInterval
	}

	//Если выставили флаг, то значение cfg.Restore у нас установлено корректно через flag.BoolVar
	//Если флаг не выставили, то cfg.Restore у нас держит значение по умолчанию через flag.BoolVar
	//Переопределяем cfg.Restore только если было выставлено значение в переменной ОС RESTORE
	//
	//$ unset RESTORE
	//$ ./server -r false # Not restore server from file
	if os.Getenv("RESTORE") != "" {
		cfg.Restore = cEnv.Restore
	}

	if cEnv.DatabaseDSN != "" {
		cfg.DatabaseDSN = cEnv.DatabaseDSN
	}

}

func GetAgentConfig() AgentConfig {
	var cfg AgentConfig
	var cEnv AgentConfigEnv

	flag.StringVar(&cfg.EndPoint, "a", serverAddressDefault, "server address")
	flag.DurationVar(&cfg.PollInterval, "p", pollIntervalDefault, "poll metrics interval")
	flag.DurationVar(&cfg.ReportInterval, "r", reportIntervalDefault, "report metric interval")
	flag.StringVar(&cfg.Key, "k", "", "key for hash function")

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

	flag.StringVar(&cfg.ServerAddress, "a", serverAddressDefault, "server address")
	flag.DurationVar(&cfg.StoreInterval, "i", storeIntervalDefault, "dump metrics to file interval")
	flag.StringVar(&cfg.StoreFile, "f", storeFileDefault, "dump file name")
	flag.BoolVar(&cfg.Restore, "r", restoreDefault, "restore metrics from file")
	flag.StringVar(&cfg.Key, "k", "", "key for hash function")
	//postgres://username:password@localhost:5432/database_name
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "Data Source Name")

	flag.Parse()
	err := env.Parse(&cEnv)
	if err != nil {
		log.Println(err)
	}

	validateServerConfig(&cfg, &cEnv)
	return cfg
}
