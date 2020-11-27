package auth

import (
	"strings"

	"github.com/dgrijalva/jwt-go"
)

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

func (c *Claims) AllowedPodsCreate() bool {
	return strings.Contains(c.Subject, "create:pods")
}

func (c *Claims) AllowedPodsRead() bool {
	return strings.Contains(c.Subject, "read:pods")
}

func (c *Claims) AllowedPodsDelete() bool {
	return strings.Contains(c.Subject, "delete:pods")
}

func (c *Claims) AllowedPeaWrite() bool {
	return strings.Contains(c.Subject, "write:pea")
}

func (c *Claims) AllowedPeaRead() bool {
	return strings.Contains(c.Subject, "read:pea")
}

func (c *Claims) AllowedPeaDelete() bool {
	return strings.Contains(c.Subject, "delete:pea")
}
