// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package handler

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/iotexproject/phoenix-gem/handler/midware"
	"github.com/iotexproject/phoenix-gem/storage"
	"go.uber.org/zap"

	"github.com/iotexproject/phoenix-gem/auth"
	"github.com/iotexproject/phoenix-gem/config"
	"github.com/iotexproject/phoenix-gem/log"
)

type StorageHandler struct {
	cfg         *config.Config
	log         *zap.Logger
	podsHandler *podsHandler
	peaHandler  *peaHandler
	cred        midware.Credential
}

func NewStorageHandler(cfg *config.Config, cred midware.Credential, provider storage.Backend) *StorageHandler {
	return &StorageHandler{
		cfg:         cfg,
		log:         log.Logger("handler"),
		podsHandler: newPodsHandler(provider),
		peaHandler:  newPeaHandler(provider),
		cred:        cred,
	}
}

func (h *StorageHandler) ServerMux(r chi.Router) http.Handler {
	r.Group(func(r chi.Router) {
		r.Use(midware.JWTTokenValid)
		r.Use(h.cred.DoCredential)
		r.Route("/pods", func(r chi.Router) {
			r.Post("/", h.podsHandler.Create)           //create bucket
			r.Delete("/{bucket}", h.podsHandler.Delete) //delete bucket
		})
		r.Route("/pea", func(r chi.Router) {
			r.Get("/{bucket}", h.peaHandler.GetObjects)        //get all objects in bucket
			r.Post("/{bucket}/*", h.peaHandler.CreateObject)   //upload object
			r.Get("/{bucket}/*", h.peaHandler.GetObject)       //download object
			r.Delete("/{bucket}/*", h.peaHandler.DeleteObject) //delete object
		})
	})
	return r
}

// for testing
func (h *StorageHandler) simpleStore(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = auth.WithStoreCtx(ctx, auth.NewStore("s3", "a", "http://localhost:9001", "AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
