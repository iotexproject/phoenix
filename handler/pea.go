package handler

import (
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/iotexproject/phoenix-gem/auth"
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
// example: curl  -H "Authorization: Bearer eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDY1NTAzMjEsImlhdCI6MTYwNjQ2MzkyMSwiaXNzIjoiMHgwNGZjZWQxYmFhNDEwODZmZTU2OGE0OGEzMzhhZjEwNGVlMTk3NzgwNDNkOThjMjI2NTU3MzRkYzg4NTkwODYxYjI2OWRlMTg3M2I3ZjhmYWM0ZGE4NjdiMjRhN2M3NDczOWZmM2Q0NmY2ZDAzYzlkYWI4YzcxMDZiYWZiOTdhODA5Iiwic3ViIjoid3JpdGU6cGVhIn0.nEFMulTTwZgJLlFE7k_lhBCKo46VCKuqkuycsG2XsVYvoYwMhxnpfpX92nCX2nQJnsiru12IW0G5QgOc_JJEVw" -d 'foobar' http://localhost:8080/pea/test11/foo.txt
func (h *peaHandler) CreateObject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := ctx.Value(auth.TokenCtxKey).(*auth.Claims)
	if !ok {
		renderJSON(w, http.StatusUnauthorized, H{"message": http.StatusText(http.StatusUnauthorized)})
		return
	}
	//check scope permission
	if !claims.AllowedPeaWrite() {
		renderJSON(w, http.StatusForbidden, H{"message": ErrorPermissionDenied.Error()})
		return
	}
	if r.Body == nil {
		renderJSON(w, http.StatusBadRequest, H{"message": ErrorBodyEmpty.Error()})
		return
	}

	bucket := chi.URLParam(r, "bucket")
	path := chi.URLParam(r, "*")
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		renderJSON(w, http.StatusBadRequest, H{"message": err.Error()})
		return
	}
	err = h.storage.PutObject(bucket, path, content)
	if err != nil {
		renderJSON(w, http.StatusBadRequest, H{"message": err.Error()})
		return
	}
	ret := H{"name": bucket, "path": path, "message": "successful"}
	renderJSON(w, http.StatusOK, ret)
}

// GetObject get pea object in storage
// example: curl -H "Authorization: Bearer eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDY1NTA2OTUsImlhdCI6MTYwNjQ2NDI5NSwiaXNzIjoiMHgwNGZjZWQxYmFhNDEwODZmZTU2OGE0OGEzMzhhZjEwNGVlMTk3NzgwNDNkOThjMjI2NTU3MzRkYzg4NTkwODYxYjI2OWRlMTg3M2I3ZjhmYWM0ZGE4NjdiMjRhN2M3NDczOWZmM2Q0NmY2ZDAzYzlkYWI4YzcxMDZiYWZiOTdhODA5Iiwic3ViIjoicmVhZDpwZWEifQ.VonbuRLKAUmHvVVdAs5Rf5d7TcPOZnWO89wMFIVLnh3jeBs77Qkg8w8_v0TMyaHA2V8OTQhOfqyWw54C6gGfyg" http://localhost:8080/pea/test11/foo.txt
func (h *peaHandler) GetObject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := ctx.Value(auth.TokenCtxKey).(*auth.Claims)
	if !ok {
		renderJSON(w, http.StatusUnauthorized, H{"message": http.StatusText(http.StatusUnauthorized)})
		return
	}
	//check scope permission
	if !claims.AllowedPeaRead() {
		renderJSON(w, http.StatusForbidden, H{"message": ErrorPermissionDenied.Error()})
		return
	}

	bucket := chi.URLParam(r, "bucket")
	path := chi.URLParam(r, "*")
	object, err := h.storage.GetObject(bucket, path)
	if err != nil {
		renderJSON(w, http.StatusBadRequest, H{"message": err.Error()})
		return
	}
	renderJSON(w, http.StatusOK, H{"message": "successful", "content": string(object.Content)})
}

// GetObjects get pea objects with bucket in storage
// example: curl -H "Authorization: Bearer eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDY1NTE0MzgsImlhdCI6MTYwNjQ2NTAzOCwiaXNzIjoiMHgwNGZjZWQxYmFhNDEwODZmZTU2OGE0OGEzMzhhZjEwNGVlMTk3NzgwNDNkOThjMjI2NTU3MzRkYzg4NTkwODYxYjI2OWRlMTg3M2I3ZjhmYWM0ZGE4NjdiMjRhN2M3NDczOWZmM2Q0NmY2ZDAzYzlkYWI4YzcxMDZiYWZiOTdhODA5Iiwic3ViIjoicmVhZDpwb2RzIn0.EdHjnkgudRBx373V6bFpA5w1dYLZmHcfkM-d7ZYqBYq3uG6W2oY0ZBpkLOPyElQe4C2r4Ual09N7AHOgVuJuLg" http://localhost:8080/pea/test11
func (h *peaHandler) GetObjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := ctx.Value(auth.TokenCtxKey).(*auth.Claims)
	if !ok {
		renderJSON(w, http.StatusUnauthorized, H{"message": http.StatusText(http.StatusUnauthorized)})
		return
	}
	//check scope permission
	if !claims.AllowedPodsRead() {
		renderJSON(w, http.StatusForbidden, H{"message": ErrorPermissionDenied.Error()})
		return
	}

	bucket := chi.URLParam(r, "bucket")
	objects, err := h.storage.ListObjects(bucket, "")
	if err != nil {
		renderJSON(w, http.StatusBadRequest, H{"message": err.Error()})
		return
	}

	list := []string{}
	for _, o := range objects {
		list = append(list, o.Path)
	}
	renderJSON(w, http.StatusOK, H{"message": "successful", "content": list})
}

// DeleteObject delete pea object in storage
// example: curl -H "Authorization: Bearer eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDY1NTEwMzksImlhdCI6MTYwNjQ2NDYzOSwiaXNzIjoiMHgwNGZjZWQxYmFhNDEwODZmZTU2OGE0OGEzMzhhZjEwNGVlMTk3NzgwNDNkOThjMjI2NTU3MzRkYzg4NTkwODYxYjI2OWRlMTg3M2I3ZjhmYWM0ZGE4NjdiMjRhN2M3NDczOWZmM2Q0NmY2ZDAzYzlkYWI4YzcxMDZiYWZiOTdhODA5Iiwic3ViIjoiZGVsZXRlOnBlYSJ9.659pN1RgCLjgGCwrSvZHpnlWVEKjj6YDJWdObCAR14p7Gr9lck-E9m7-U3stRm10jYAjUVQFUUQtJNzWLxv3mQ" -X DELETE http://localhost:8080/pea/test11/foo.txt
func (h *peaHandler) DeleteObject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := ctx.Value(auth.TokenCtxKey).(*auth.Claims)
	if !ok {
		renderJSON(w, http.StatusUnauthorized, H{"message": http.StatusText(http.StatusUnauthorized)})
		return
	}
	//check scope permission
	if !claims.AllowedPeaDelete() {
		renderJSON(w, http.StatusForbidden, H{"message": ErrorPermissionDenied.Error()})
		return
	}

	bucket := chi.URLParam(r, "bucket")
	path := chi.URLParam(r, "*")
	err := h.storage.DeleteObject(bucket, path)
	if err != nil {
		renderJSON(w, http.StatusBadRequest, H{"message": err.Error()})
		return
	}
	renderJSON(w, http.StatusOK, H{"message": "successful", "name": bucket, "path": path})
}
