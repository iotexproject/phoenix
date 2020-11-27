package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/iotexproject/go-pkgs/crypto"
	"github.com/iotexproject/phoenix-gem/auth"
)

// JWTTokenValid operation middleware
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

		claim := &auth.Claims{}

		token, err := jwt.ParseWithClaims(jwtString, claim, func(token *jwt.Token) (interface{}, error) {
			keyHex := claim.Issuer
			if keyHex[:2] == "0x" || keyHex[:2] == "0X" {
				keyHex = keyHex[2:]
			}
			key, err := crypto.HexStringToPublicKey(keyHex)
			if err != nil {
				return nil, err
			}
			return key.EcdsaPublicKey(), nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if !token.Valid || len(token.Header) < 2 {
			// should not happen with a success parsing, check anyway
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, auth.TokenCtxKey, claim)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
