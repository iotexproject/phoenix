// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

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
	return strings.Contains(c.Scope, "create:pods")
}

func (c *Claims) AllowedPodsRead() bool {
	return strings.Contains(c.Scope, "read:pods")
}

func (c *Claims) AllowedPodsDelete() bool {
	return strings.Contains(c.Scope, "delete:pods")
}

func (c *Claims) AllowedPeaWrite() bool {
	return strings.Contains(c.Scope, "write:pea")
}

func (c *Claims) AllowedPeaRead() bool {
	return strings.Contains(c.Scope, "read:pea")
}

func (c *Claims) AllowedPeaDelete() bool {
	return strings.Contains(c.Scope, "delete:pea")
}
