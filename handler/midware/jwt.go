// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package midware

import (
	"context"
	"net/http"
	"strings"

	"github.com/iotexproject/iotex-antenna-go/v2/jwt"

	"github.com/iotexproject/phoenix-gem/auth"
)

// JWTTokenValid operation midware
func JWTTokenValid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from authorization header.
		jwtString := ""
		bearer := r.Header.Get("Authorization")
		if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
			jwtString = bearer[7:]
		}

		if jwtString == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		jwtoken, err := jwt.VerifyJWT(jwtString)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), auth.TokenCtxKey, &auth.Claims{JWT: jwtoken})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
