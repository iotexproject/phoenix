// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package auth

import (
	"testing"

	"github.com/iotexproject/iotex-antenna-go/v2/jwt"
	"github.com/stretchr/testify/require"
)

func TestClaim(t *testing.T) {
	r := require.New(t)

	c := Claims{JWT: &jwt.JWT{}}

	c.Scope = "Create: 0ZTpwb2RzIiwiZXhwI"
	r.True(c.AllowCreate())
	c.Scope = "Read: wNzY2OTI0OSwia"
	r.True(c.AllowRead())
	c.Scope = "Delete: IiOiJodHRwOi8"
	r.True(c.AllowDelete())
	c.Scope = "Update: jb21lLzEyMzQif"
	r.True(c.AllowWrite())
	c.Subject = "Bucket: NiIsInR5cCI6Ikp"
	r.True(c.IsBucket())
	c.Subject = "Object: nh74PpkJJibjA"
	r.True(c.IsObject())
}
