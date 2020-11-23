package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/iotexproject/phoenix-gem/json"
	"github.com/iotexproject/phoenix-gem/log"
	"github.com/iotexproject/phoenix-gem/storage"
	"go.uber.org/zap"
)

type podsHandler struct {
	log     *zap.Logger
	storage storage.Backend
}

type podObject struct {
	Name string `json:"name"`
}

func newPodsHandler(provider storage.Backend) *podsHandler {
	return &podsHandler{
		log:     log.Logger("pods"),
		storage: provider,
	}
}

// Create create pod in storage
// example: curl -H "Content-type: application/json" -d '{ "name": "test11"}' http://localhost:8080/pods
func (h *podsHandler) Create(w http.ResponseWriter, r *http.Request) {
	item := &podObject{}
	if err := decodeJSON(r, item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

// Delete delete pod in storage
// example: curl -X DELETE http://localhost:8080/pods/test11
func (h *podsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	bucket := chi.URLParam(r, "bucket")
	err := h.storage.DeleteBucket(bucket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
