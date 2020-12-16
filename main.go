// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"

	"github.com/iotexproject/phoenix/config"
	"github.com/iotexproject/phoenix/log"
	"github.com/iotexproject/phoenix/server"
)

const (
	ConfigPath = "ConfigPath"
)

func main() {
	configPath := os.Getenv(ConfigPath)
	if configPath == "" {
		configPath = "/var/data/config.yaml"
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
	go func() {
		if err = srv.Start(); err != nil {
			log.L().Fatal("server start:", zap.Error(err))
		}
	}()
	handleShutdown(srv)
}

type Stopper interface {
	Stop(context.Context) error
}

func handleShutdown(service ...Stopper) {
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, os.Kill)

	// wait INT or KILL
	<-stop
	log.L().Info("shutting down ...")
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	for _, s := range service {
		if err := s.Stop(ctx); err != nil {
			log.L().Error("shutdown", zap.Error(err))
		}
	}
	return
}
