// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"

	"github.com/iotexproject/phoenix/auth"
	"github.com/iotexproject/phoenix/config"
	"github.com/iotexproject/phoenix/db"
	"github.com/iotexproject/phoenix/handler"
	"github.com/iotexproject/phoenix/log"
	"github.com/iotexproject/phoenix/worker"
)

// Server struct
type Server struct {
	*http.Server
	cfg    *config.Config
	log    *zap.Logger
	userDB db.KVStore
	//queue  diskqueue.Interface
	worker *worker.Worker
}

// New return new Server instance
func New(cfg *config.Config) *Server {
	srv := &Server{
		cfg: cfg,
		log: log.Logger("server"),
	}
	return srv
}

// Start start the server
func (srv *Server) Start() error {
	r := chi.NewRouter()
	// middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	srv.log.Info("RateLimit",
		zap.Bool("Enable", srv.cfg.Server.RateLimit.Enable),
		zap.Strings("LimitByKey", srv.cfg.Server.RateLimit.LimitByKey),
		zap.Int("RequestLimit", srv.cfg.Server.RateLimit.RequestLimit),
		zap.Int("WindowLength", srv.cfg.Server.RateLimit.WindowLength),
	)
	srv.log.Info("CORS",
		zap.Bool("Enable", srv.cfg.Server.Cors.Enable),
		zap.Strings("AllowedOrigins", srv.cfg.Server.Cors.AllowedOrigins),
		zap.Strings("AllowedMethods", srv.cfg.Server.Cors.AllowedMethods),
		zap.Strings("AllowedHeaders", srv.cfg.Server.Cors.AllowedHeaders),
	)
	if srv.cfg.Server.Cors.Enable {
		r.Use(cors.Handler(cors.Options{
			//Debug:            true,
			AllowedOrigins:   srv.cfg.Server.Cors.AllowedOrigins,
			AllowedMethods:   srv.cfg.Server.Cors.AllowedMethods,
			AllowedHeaders:   srv.cfg.Server.Cors.AllowedHeaders,
			AllowCredentials: false,
		}))
	}

	var err error
	// open db for user's storage endpoint
	srv.userDB = db.NewBoltDB(srv.cfg.Server.DBPath)
	if err = srv.userDB.Start(context.Background()); err != nil {
		srv.log.Error("start db", zap.Error(err))
		return err
	}

	srv.worker = worker.New(5).Start(context.Background(), func(job worker.Job) { fmt.Printf("%s", job.Data) })
	endpoint := fmt.Sprintf(":%s", srv.cfg.Server.Port)
	h := handler.NewStorageHandler(srv.cfg, auth.NewCredential(srv.userDB))
	srv.Server = &http.Server{
		Handler: h.ServerMux(r),
		Addr:    endpoint,
	}
	srv.log.Info("starting server", zap.String("endpoint", endpoint))
	return srv.ListenAndServe()
}

func (srv *Server) Stop(ctx context.Context) error {
	if err := srv.Shutdown(ctx); err != nil {
		srv.log.Error("shutdown server", zap.Error(err))
	}
	return srv.userDB.Stop(ctx)
}
