package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/iotexproject/phoenix-gem/config"
	"github.com/iotexproject/phoenix-gem/handler"
	"github.com/iotexproject/phoenix-gem/log"
	"github.com/rs/cors"
	"go.uber.org/zap"
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
	srv.log.Debug("enter server")
	r := chi.NewRouter()
	// Basic CORS
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.AllowAll().Handler)

	endpoint := fmt.Sprintf(":%s", srv.cfg.Server.Port)
	h := handler.NewStorageHandler(srv.cfg)
	s := &http.Server{
		Handler: h.ServerMux(r),
		Addr:    endpoint,
	}
	srv.log.Info("starting server", zap.String("endpoint", endpoint))
	return s.ListenAndServe()
}
