// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package auth

import (
	"strings"

	"github.com/iotexproject/iotex-antenna-go/v2/jwt"
)

// const
const (
	Bucket = "Bucket"
	Object = "Object"
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
	*jwt.JWT
}

func (c *Claims) AllowCreate() bool {
	return strings.Contains(c.Scope, jwt.CREATE)
}

func (c *Claims) AllowRead() bool {
	return strings.Contains(c.Scope, jwt.READ)
}

func (c *Claims) AllowDelete() bool {
	return strings.Contains(c.Scope, jwt.DELETE)
}

func (c *Claims) AllowWrite() bool {
	return strings.Contains(c.Scope, jwt.UPDATE)
}

func (c *Claims) IsBucket() bool {
	return strings.Contains(c.Subject, Bucket)
}

func (c *Claims) IsObject() bool {
	return strings.Contains(c.Subject, Object)
}
