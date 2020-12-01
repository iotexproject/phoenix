// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"go.uber.org/zap"

	"github.com/iotexproject/phoenix-gem/config"
	"github.com/iotexproject/phoenix-gem/handler"
	"github.com/iotexproject/phoenix-gem/log"
	"github.com/iotexproject/phoenix-gem/storage"
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

	endpoint := fmt.Sprintf(":%s", srv.cfg.Server.Port)
	provider, err := srv.getProvider()
	if err != nil {
		return err
	}
	h := handler.NewStorageHandler(srv.cfg, provider)

	s := &http.Server{
		Handler: h.ServerMux(r),
		Addr:    endpoint,
	}
	srv.log.Info("starting server", zap.String("endpoint", endpoint))
	return s.ListenAndServe()
}

func (srv *Server) getProvider() (storage.Backend, error) {
	var provider storage.Backend
	var err error
	switch srv.cfg.Storage.Provider {
	case "s3":
		scr := credentials.NewStaticCredentials(
			srv.cfg.S3.AccessKey,
			srv.cfg.S3.SecretKey,
			"")
		provider = storage.NewAmazonS3BackendWithCredentials("", srv.cfg.S3.Region, srv.cfg.S3.EndPoint, "", scr)
	default:
		err = fmt.Errorf("storage provider `%s` not supported", srv.cfg.Storage.Provider)
	}
	return provider, err
}
