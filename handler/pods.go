// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/iotexproject/phoenix-gem/auth"
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
// example: curl -H "Authorization: Bearer eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDY1NDg0MzAsImlhdCI6MTYwNjQ2MjAzMCwiaXNzIjoiMHgwNGZjZWQxYmFhNDEwODZmZTU2OGE0OGEzMzhhZjEwNGVlMTk3NzgwNDNkOThjMjI2NTU3MzRkYzg4NTkwODYxYjI2OWRlMTg3M2I3ZjhmYWM0ZGE4NjdiMjRhN2M3NDczOWZmM2Q0NmY2ZDAzYzlkYWI4YzcxMDZiYWZiOTdhODA5Iiwic3ViIjoiY3JlYXRlOnBvZHMifQ.BeQs3s6rx-x9O-JURRwnbUzHIcvDSF0TSqJ5GsBZPLnZPE0s3rQvgA8wmZlTNNJzvfL3hUZSfwY3-Lg-vvaxfg" -H "Content-type: application/json" -d '{ "name": "test10"}' http://localhost:8080/pods
func (h *podsHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := ctx.Value(auth.TokenCtxKey).(*auth.Claims)
	if !ok {
		renderJSON(w, http.StatusUnauthorized, H{"message": http.StatusText(http.StatusUnauthorized)})
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

	_, err := h.storage.CreateBucket(item.Name)
	if err != nil {
		renderJSON(w, http.StatusBadRequest, H{"message": err.Error()})
		return
	}

	ret := H{"name": item.Name, "message": "successful"}
	renderJSON(w, http.StatusOK, ret)
}

// Delete delete pod in storage
// example: curl  -H "Authorization: Bearer eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDY1NDk3OTcsImlhdCI6MTYwNjQ2MzM5NywiaXNzIjoiMHgwNGZjZWQxYmFhNDEwODZmZTU2OGE0OGEzMzhhZjEwNGVlMTk3NzgwNDNkOThjMjI2NTU3MzRkYzg4NTkwODYxYjI2OWRlMTg3M2I3ZjhmYWM0ZGE4NjdiMjRhN2M3NDczOWZmM2Q0NmY2ZDAzYzlkYWI4YzcxMDZiYWZiOTdhODA5Iiwic3ViIjoiZGVsZXRlOnBvZHMifQ.LVDnJPCql_0dlJj_mDqSjFZL9V46lWpp37aY_GCem3TR3625ZrZ-6mEEcNN4N94RW02fRm0AdwDbR2Iz0BJL_Q" -X DELETE http://localhost:8080/pods/test111
func (h *podsHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
	bucket := chi.URLParam(r, "bucket")
	err := h.storage.DeleteBucket(bucket)
	if err != nil {
		renderJSON(w, http.StatusBadRequest, H{"message": err.Error()})
		return
	}
	ret := H{"name": bucket, "message": "successful"}
	renderJSON(w, http.StatusOK, ret)
}
