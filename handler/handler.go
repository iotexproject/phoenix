package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/iotexproject/phoenix-gem/models"
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
	cfg     *config.Config
	log     *zap.Logger
	storage storage.Backend
}

func decodeJSON(r *http.Request, v interface{}) error {
	defer io.Copy(ioutil.Discard, r.Body)
	return json.NewDecoder(r.Body).Decode(v)
}

func NewStorageHandler(cfg *config.Config, provider storage.Backend) *StorageHandler {

	return &StorageHandler{
		cfg:     cfg,
		log:     log.Logger("handler"),
		storage: provider,
	}
}

func (h *StorageHandler) ServerMux(r chi.Router) http.Handler {
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTTokenValid)
		r.Put("/putobject/{pod}/*", h.putObject)
		r.Post("/createobject", h.createObject)
	})
	return r
}

/*
curl -H "Content-type: application/json" -d '{ "name": "test11"}' http://localhost:8080/createobject

*/
func (h *StorageHandler) createObject(w http.ResponseWriter, r *http.Request) {
	item := &models.PutObject{}
	if err := decodeJSON(r, item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.log.Debug("receive createObject", zap.Any("req", fmt.Sprintf("%v", r)), zap.String("name", item.Name))

	object, err := h.storage.CreateBucket(item.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	ret := object

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ret)

}

/*
curl -d 'foobar' -X PUT -H 'Content-Type: text/plain' http://localhost:8080/putobject/test11/foo.txt

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

	err = h.storage.PutObject(pod, pea, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	ret := H{"url": "http://test.com/"}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ret)

}
