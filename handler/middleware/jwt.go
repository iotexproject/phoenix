package middleware

import (
	"fmt"
	"net/http"
)

// JWTTokenValid operation middleware
func JWTTokenValid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var err error

		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid format for parameter identifier: %s", err), http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
