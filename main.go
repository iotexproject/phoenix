package main

import (
	"fmt"
	"os"

	"github.com/iotexproject/phoenix-gem/config"
	"github.com/iotexproject/phoenix-gem/log"
	"github.com/iotexproject/phoenix-gem/server"
	"go.uber.org/zap"
)

const (
	ConfigPath = "ConfigPath"
)

func main() {
	configPath := os.Getenv(ConfigPath)
	if configPath == "" {
		configPath = "config.yaml"
	}
	cfg, err := config.New(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to parse config: %v\n", err)
		os.Exit(1)
	}

	if err := log.InitLoggers(cfg.Log, cfg.SubLogs); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to init logger: %v\n", err)
		os.Exit(1)
	}
	log.L().Info("init logger")

	srv := server.New(cfg)
	if err = srv.Start(); err != nil {
		log.L().Fatal("server start:", zap.Error(err))
	}
}
