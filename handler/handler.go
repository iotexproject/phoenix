package handler

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/iotexproject/phoenix-gem/storage"
	"go.uber.org/zap"

	"github.com/iotexproject/phoenix-gem/config"
	"github.com/iotexproject/phoenix-gem/handler/middleware"
	"github.com/iotexproject/phoenix-gem/json"
	"github.com/iotexproject/phoenix-gem/log"
)

// H is a shortcut for map[string]interface{}
type H map[string]interface{}

type StorageHandler struct {
	cfg         *config.Config
	log         *zap.Logger
	podsHandler *podsHandler
	peaHandler  *peaHandler
}

func decodeJSON(r *http.Request, v interface{}) error {
	defer io.Copy(ioutil.Discard, r.Body)
	return json.NewDecoder(r.Body).Decode(v)
}

func NewStorageHandler(cfg *config.Config, provider storage.Backend) *StorageHandler {
	return &StorageHandler{
		cfg:         cfg,
		log:         log.Logger("handler"),
		podsHandler: newPodsHandler(provider),
		peaHandler:  newPeaHandler(provider),
	}
}

func (h *StorageHandler) ServerMux(r chi.Router) http.Handler {
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTTokenValid)
		r.Route("/pods", func(r chi.Router) {
			r.Post("/", h.podsHandler.Create)           //create bucket
			r.Delete("/{bucket}", h.podsHandler.Delete) //delete bucket
		})
		r.Route("/pea", func(r chi.Router) {
			r.Get("/{bucket}", h.peaHandler.GetObjects)        //get objects in bucket
			r.Post("/{bucket}/*", h.peaHandler.CreateObject)   //upload object
			r.Get("/{bucket}/*", h.peaHandler.GetObject)       //upload object
			r.Delete("/{bucket}/*", h.peaHandler.DeleteObject) //delete object
		})
	})
	return r
}
