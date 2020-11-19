package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/iotexproject/phoenix-gem/config"
	"github.com/rs/cors"
)

// Server struct
type Server struct {
	cfg *config.Config
}

// New return new Server instance
func New(cfg *config.Config) *Server {
	srv := &Server{cfg: cfg}
	return srv
}

// Start start server
func (srv *Server) Start() error {
	r := chi.NewRouter()
	// Basic CORS
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.AllowAll().Handler)

	endpoint := fmt.Sprintf(":%s", srv.cfg.Server.Port)
	s := &http.Server{
		Handler: http.DefaultServeMux,
		Addr:    endpoint,
	}

	log.Println("listen at ", endpoint)
	return s.ListenAndServe()
}
