// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package handler

import (
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/iotexproject/phoenix-gem/auth"
	"github.com/iotexproject/phoenix-gem/handler/midware"
	"github.com/iotexproject/phoenix-gem/storage"
	"go.uber.org/zap"

	"github.com/iotexproject/phoenix-gem/config"
	"github.com/iotexproject/phoenix-gem/log"
)

type StorageHandler struct {
	cfg  *config.Config
	log  *zap.Logger
	cred midware.Credential
}

func NewStorageHandler(cfg *config.Config, cred midware.Credential) *StorageHandler {
	return &StorageHandler{
		cfg:  cfg,
		log:  log.Logger("handler"),
		cred: cred,
	}
}

func (h *StorageHandler) ServerMux(r chi.Router) http.Handler {
	r.Group(func(r chi.Router) {
		r.Use(midware.JWTTokenValid)
		r.Use(h.cred.DoCredential)
		r.Route("/pods", func(r chi.Router) {
			r.Post("/", h.CreateBucket)           //create bucket
			r.Delete("/{bucket}", h.DeleteBucket) //delete bucket
		})
		r.Route("/pea", func(r chi.Router) {
			r.Get("/{bucket}", h.GetObjects)        //get all objects in bucket
			r.Post("/{bucket}/*", h.CreateObject)   //upload object
			r.Get("/{bucket}/*", h.GetObject)       //download object
			r.Delete("/{bucket}/*", h.DeleteObject) //delete object
		})
	})
	return r
}

// CreateBucket create pod in storage
// example: curl -H "Authorization: Bearer jwttoken" -H "Content-type: application/json" -d '{ "name": "test10"}' http://localhost:8080/pods
func (h *StorageHandler) CreateBucket(w http.ResponseWriter, r *http.Request) {
	claims, storage, statusCode := h.createBackendForRequest(r)
	if statusCode != http.StatusOK {
		renderJSON(w, statusCode, http.StatusText(statusCode))
		return
	}
	//check scope permission
	if !claims.AllowCreate() {
		renderJSON(w, http.StatusForbidden, H{"message": ErrorPermissionDenied.Error()})
		return
	}
	item := &podObject{}
	if err := decodeJSON(r, item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := storage.CreateBucket(item.Name)
	if err != nil {
		renderJSON(w, http.StatusBadRequest, H{"message": err.Error()})
		return
	}

	ret := H{"name": item.Name, "message": "successful"}
	renderJSON(w, http.StatusOK, ret)
}

// DeleteBucket delete pod in storage
// example: curl  -H "Authorization: Bearer jwttoken" -X DELETE http://localhost:8080/pods/test111
func (h *StorageHandler) DeleteBucket(w http.ResponseWriter, r *http.Request) {
	claims, storage, statusCode := h.createBackendForRequest(r)
	if statusCode != http.StatusOK {
		renderJSON(w, statusCode, http.StatusText(statusCode))
		return
	}
	//check scope permission
	if !claims.AllowDelete() {
		renderJSON(w, http.StatusForbidden, H{"message": ErrorPermissionDenied.Error()})
		return
	}

	bucket := chi.URLParam(r, "bucket")
	err := storage.DeleteBucket(bucket)
	if err != nil {
		renderJSON(w, http.StatusBadRequest, H{"message": err.Error()})
		return
	}
	ret := H{"name": bucket, "message": "successful"}
	renderJSON(w, http.StatusOK, ret)
}

// CreateObject create pod in storage
// example: curl  -H "Authorization: Bearer jwttoken" -d 'foobar' http://localhost:8080/pea/test11/foo.txt
func (h *StorageHandler) CreateObject(w http.ResponseWriter, r *http.Request) {
	claims, storage, statusCode := h.createBackendForRequest(r)
	if statusCode != http.StatusOK {
		renderJSON(w, statusCode, http.StatusText(statusCode))
		return
	}
	//check scope permission
	if !claims.AllowWrite() {
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
	err = storage.PutObject(bucket, path, content)
	if err != nil {
		renderJSON(w, http.StatusBadRequest, H{"message": err.Error()})
		return
	}
	ret := H{"name": bucket, "path": path, "message": "successful"}
	renderJSON(w, http.StatusOK, ret)
}

// GetObject get pea object in storage
// example: curl -H "Authorization: Bearer jwttoken" http://localhost:8080/pea/test11/foo.txt
func (h *StorageHandler) GetObject(w http.ResponseWriter, r *http.Request) {
	claims, storage, statusCode := h.createBackendForRequest(r)
	if statusCode != http.StatusOK {
		renderJSON(w, statusCode, http.StatusText(statusCode))
		return
	}
	//check scope permission
	if !claims.AllowRead() {
		renderJSON(w, http.StatusForbidden, H{"message": ErrorPermissionDenied.Error()})
		return
	}

	bucket := chi.URLParam(r, "bucket")
	path := chi.URLParam(r, "*")
	object, err := storage.GetObject(bucket, path)
	if err != nil {
		renderJSON(w, http.StatusBadRequest, H{"message": err.Error()})
		return
	}
	renderJSON(w, http.StatusOK, H{"message": "successful", "content": string(object.Content)})
}

// GetObjects get pea objects with bucket in storage
// example: curl -H "Authorization: Bearer jwttoken" http://localhost:8080/pea/test11
func (h *StorageHandler) GetObjects(w http.ResponseWriter, r *http.Request) {
	claims, storage, statusCode := h.createBackendForRequest(r)
	if statusCode != http.StatusOK {
		renderJSON(w, statusCode, http.StatusText(statusCode))
		return
	}
	//check scope permission
	if !claims.AllowRead() {
		renderJSON(w, http.StatusForbidden, H{"message": ErrorPermissionDenied.Error()})
		return
	}

	bucket := chi.URLParam(r, "bucket")
	objects, err := storage.ListObjects(bucket, "")
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
// example: curl -H "Authorization: Bearer jwttoken" -X DELETE http://localhost:8080/pea/test11/foo.txt
func (h *StorageHandler) DeleteObject(w http.ResponseWriter, r *http.Request) {
	claims, storage, statusCode := h.createBackendForRequest(r)
	if statusCode != http.StatusOK {
		renderJSON(w, statusCode, http.StatusText(statusCode))
		return
	}
	//check scope permission
	if !claims.AllowDelete() {
		renderJSON(w, http.StatusForbidden, H{"message": ErrorPermissionDenied.Error()})
		return
	}

	bucket := chi.URLParam(r, "bucket")
	path := chi.URLParam(r, "*")
	err := storage.DeleteObject(bucket, path)
	if err != nil {
		renderJSON(w, http.StatusBadRequest, H{"message": err.Error()})
		return
	}
	renderJSON(w, http.StatusOK, H{"message": "successful", "name": bucket, "path": path})
}

func (h *StorageHandler) createBackendForRequest(r *http.Request) (claims *auth.Claims, backend storage.Backend, statusCode int) {
	ctx := r.Context()
	claims, ok := ctx.Value(auth.TokenCtxKey).(*auth.Claims)
	if !ok {
		statusCode = http.StatusBadRequest
		return
	}

	store, ok := auth.GetStoreCtx(ctx)
	if !ok {
		statusCode = http.StatusBadRequest
		return
	}
	backend, err := storage.NewStorage(store)
	if err != nil {
		h.log.Error("failed to new storage", zap.Error(err))
		statusCode = http.StatusServiceUnavailable
		return
	}

	statusCode = http.StatusOK
	return
}
