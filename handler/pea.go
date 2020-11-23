package handler

import (
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/iotexproject/phoenix-gem/json"
	"github.com/iotexproject/phoenix-gem/log"
	"github.com/iotexproject/phoenix-gem/storage"
	"go.uber.org/zap"
)

type peaHandler struct {
	log     *zap.Logger
	storage storage.Backend
}

type peaObject struct {
	Name string `json:"name"`
}

func newPeaHandler(provider storage.Backend) *peaHandler {
	return &peaHandler{
		log:     log.Logger("pea"),
		storage: provider,
	}
}

// CreateObject create pod in storage
// example: curl -d 'foobar' http://localhost:8080/pea/test11/foo.txt
func (h *peaHandler) CreateObject(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Body must be set", http.StatusBadRequest)
		return
	}

	bucket := chi.URLParam(r, "bucket")
	path := chi.URLParam(r, "*")
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.storage.PutObject(bucket, path, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetObject get pea object in storage
// example: curl http://localhost:8080/pea/test11/foo.txt
func (h *peaHandler) GetObject(w http.ResponseWriter, r *http.Request) {
	bucket := chi.URLParam(r, "bucket")
	path := chi.URLParam(r, "*")
	object, err := h.storage.GetObject(bucket, path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(object.Content)
}

// GetObjects get pea objects with bucket in storage
// example: curl http://localhost:8080/pea/test11
func (h *peaHandler) GetObjects(w http.ResponseWriter, r *http.Request) {
	bucket := chi.URLParam(r, "bucket")
	objects, err := h.storage.ListObjects(bucket, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(objects)
}

// DeleteObject delete pea object in storage
// example: curl -X DELETE http://localhost:8080/pea/test11/foo.txt
func (h *peaHandler) DeleteObject(w http.ResponseWriter, r *http.Request) {
	bucket := chi.URLParam(r, "bucket")
	path := chi.URLParam(r, "*")
	err := h.storage.DeleteObject(bucket, path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
