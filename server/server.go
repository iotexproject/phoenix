package server

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/iotexproject/phoenix-gem/config"
	"github.com/iotexproject/phoenix-gem/handler"
	"github.com/iotexproject/phoenix-gem/log"
	"github.com/iotexproject/phoenix-gem/storage"
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
