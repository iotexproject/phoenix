package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/iotexproject/phoenix-gem/config"
	"github.com/iotexproject/phoenix-gem/handler/middleware"
	"github.com/iotexproject/phoenix-gem/json"
	"github.com/iotexproject/phoenix-gem/log"
	"go.uber.org/zap"
)

// H is a shortcut for map[string]interface{}
type H map[string]interface{}

type StorageHandler struct {
	cfg *config.Config
	log *zap.Logger
}

func NewStorageHandler(cfg *config.Config) *StorageHandler {
	return &StorageHandler{
		cfg: cfg,
		log: log.Logger("handler"),
	}
}

func (h *StorageHandler) ServerMux(r chi.Router) http.Handler {
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTTokenValid)
		r.Put("/putobject/{pod}/*", h.putObject)
	})

	return r
}

/*
curl -d 'foobar' -X PUT -H 'Content-Type: text/plain' http://localhost:8080/putobject/aaa/foo.txt

*/
func (h *StorageHandler) putObject(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Body must be set", http.StatusBadRequest)
		return
	}

	pod := chi.URLParam(r, "pod")
	pea := chi.URLParam(r, "*")
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.log.Debug("receive putObject", zap.Any("req", fmt.Sprintf("%v", r)), zap.String("pod", pod), zap.String("pea", pea), zap.ByteString("content", content))
	w.Header().Set("Content-Type", "application/json")

	ret := H{"url": "http://test.com/"}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ret)

}
