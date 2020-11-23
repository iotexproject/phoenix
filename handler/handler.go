package handler

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/iotexproject/phoenix-gem/auth"
	"github.com/iotexproject/phoenix-gem/storage"
	"go.uber.org/zap"

	"github.com/iotexproject/phoenix-gem/config"
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
		r.Use(h.JWTTokenValid)
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

// JWTTokenValid operation middleware
func (h *StorageHandler) JWTTokenValid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from authorization header.
		tokenString := ""
		bearer := r.Header.Get("Authorization")
		if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
			tokenString = bearer[7:]
		}

		if tokenString == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}

			return []byte(h.cfg.Server.AuthSecret), nil
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return

		}
		ctx := r.Context()
		if claims, ok := token.Claims.(*auth.Claims); ok && token.Valid {
			ctx = context.WithValue(ctx, auth.TokenCtxKey, claims)
		} else {
			http.Error(w, fmt.Sprintf("Invalid format for parameter identifier: %s", err), http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
