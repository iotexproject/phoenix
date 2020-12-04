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
	"go.uber.org/zap"

	"github.com/iotexproject/phoenix-gem/auth"
	"github.com/iotexproject/phoenix-gem/log"
	"github.com/iotexproject/phoenix-gem/storage"
)

type peaHandler struct {
	log *zap.Logger
}

type peaObject struct {
	Name string `json:"name"`
}

func newPeaHandler() *peaHandler {
	return &peaHandler{
		log: log.Logger("pea"),
	}
}

// CreateObject create pod in storage
// example: curl  -H "Authorization: Bearer jwttoken" -d 'foobar' http://localhost:8080/pea/test11/foo.txt
func (h *peaHandler) CreateObject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := ctx.Value(auth.TokenCtxKey).(*auth.Claims)
	if !ok {
		renderJSON(w, http.StatusUnauthorized, H{"message": http.StatusText(http.StatusUnauthorized)})
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

	store, ok := auth.GetStoreCtx(ctx)
	if !ok {
		renderJSON(w, http.StatusBadRequest, H{"message": ErrorStoreCtx.Error()})
		return
	}
	storage, err := storage.NewStorage(store)
	if err != nil {
		renderJSON(w, http.StatusBadRequest, H{"message": err.Error()})
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
func (h *peaHandler) GetObject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := ctx.Value(auth.TokenCtxKey).(*auth.Claims)
	if !ok {
		renderJSON(w, http.StatusUnauthorized, H{"message": http.StatusText(http.StatusUnauthorized)})
		return
	}
	//check scope permission
	if !claims.AllowRead() {
		renderJSON(w, http.StatusForbidden, H{"message": ErrorPermissionDenied.Error()})
		return
	}

	store, ok := auth.GetStoreCtx(ctx)
	if !ok {
		renderJSON(w, http.StatusBadRequest, H{"message": ErrorStoreCtx.Error()})
		return
	}
	storage, err := storage.NewStorage(store)
	if err != nil {
		renderJSON(w, http.StatusBadRequest, H{"message": err.Error()})
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
func (h *peaHandler) GetObjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := ctx.Value(auth.TokenCtxKey).(*auth.Claims)
	if !ok {
		renderJSON(w, http.StatusUnauthorized, H{"message": http.StatusText(http.StatusUnauthorized)})
		return
	}
	//check scope permission
	if !claims.AllowRead() {
		renderJSON(w, http.StatusForbidden, H{"message": ErrorPermissionDenied.Error()})
		return
	}

	store, ok := auth.GetStoreCtx(ctx)
	if !ok {
		renderJSON(w, http.StatusBadRequest, H{"message": ErrorStoreCtx.Error()})
		return
	}
	storage, err := storage.NewStorage(store)
	if err != nil {
		renderJSON(w, http.StatusBadRequest, H{"message": err.Error()})
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
func (h *peaHandler) DeleteObject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := ctx.Value(auth.TokenCtxKey).(*auth.Claims)
	if !ok {
		renderJSON(w, http.StatusUnauthorized, H{"message": http.StatusText(http.StatusUnauthorized)})
		return
	}
	//check scope permission
	if !claims.AllowDelete() {
		renderJSON(w, http.StatusForbidden, H{"message": ErrorPermissionDenied.Error()})
		return
	}

	store, ok := auth.GetStoreCtx(ctx)
	if !ok {
		renderJSON(w, http.StatusBadRequest, H{"message": ErrorStoreCtx.Error()})
		return
	}
	storage, err := storage.NewStorage(store)
	if err != nil {
		renderJSON(w, http.StatusBadRequest, H{"message": err.Error()})
		return
	}

	bucket := chi.URLParam(r, "bucket")
	path := chi.URLParam(r, "*")
	err = storage.DeleteObject(bucket, path)
	if err != nil {
		renderJSON(w, http.StatusBadRequest, H{"message": err.Error()})
		return
	}
	renderJSON(w, http.StatusOK, H{"message": "successful", "name": bucket, "path": path})
}
