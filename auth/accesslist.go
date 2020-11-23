package auth

import "github.com/dgrijalva/jwt-go"

// Context keys
var (
	TokenCtxKey = &contextKey{"Token"}
	ErrorCtxKey = &contextKey{"Error"}
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "auth context value " + k.name
}

type Claims struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}
