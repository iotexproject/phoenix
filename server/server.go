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
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"go.uber.org/zap"

	"github.com/iotexproject/phoenix-gem/config"
	"github.com/iotexproject/phoenix-gem/db"
	"github.com/iotexproject/phoenix-gem/handler"
	"github.com/iotexproject/phoenix-gem/handler/midware"
	"github.com/iotexproject/phoenix-gem/log"
)

// Server struct
type Server struct {
	cfg *config.Config
	log *zap.Logger
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
	if srv.cfg.Server.RateLimit.Enable && srv.cfg.Server.RateLimit.RequestLimit > 0 && srv.cfg.Server.RateLimit.WindowLength > 0 {
		srv.log.Info("RateLimit enable",
			zap.Int("RequestLimit", srv.cfg.Server.RateLimit.RequestLimit),
			zap.Int("WindowLength", srv.cfg.Server.RateLimit.WindowLength),
		)
		r.Use(httprate.LimitByIP(srv.cfg.Server.RateLimit.RequestLimit, time.Duration(srv.cfg.Server.RateLimit.WindowLength)*time.Second))
	}
	if srv.cfg.Server.Cors.Enable {
		srv.log.Info("CORS enable",
			zap.Strings("AllowedOrigins", srv.cfg.Server.Cors.AllowedOrigins),
			zap.Strings("AllowedMethods", srv.cfg.Server.Cors.AllowedMethods),
			zap.Strings("AllowedHeaders", srv.cfg.Server.Cors.AllowedHeaders),
		)
		r.Use(cors.Handler(cors.Options{
			//Debug:            true,
			AllowedOrigins:   srv.cfg.Server.Cors.AllowedOrigins,
			AllowedMethods:   srv.cfg.Server.Cors.AllowedMethods,
			AllowedHeaders:   srv.cfg.Server.Cors.AllowedHeaders,
			AllowCredentials: false,
		}))
	}

	// open db for user's storage endpoint
	user := db.NewBoltDB(srv.cfg.Server.DBPath)
	ctx := context.Background()
	if err := user.Start(ctx); err != nil {
		return err
	}

	endpoint := fmt.Sprintf(":%s", srv.cfg.Server.Port)

	h := handler.NewStorageHandler(srv.cfg, midware.NewCredential(user))

	s := &http.Server{
		Handler: h.ServerMux(r),
		Addr:    endpoint,
	}
	srv.log.Info("starting server", zap.String("endpoint", endpoint))
	return s.ListenAndServe()
}
