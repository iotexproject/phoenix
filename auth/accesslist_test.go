// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package auth

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/iotexproject/go-pkgs/crypto"
	"github.com/stretchr/testify/require"
)

// SignJWT creates a JWT
func SignJWT(issue, expire int64, subject, scope string, key crypto.PrivateKey) (string, error) {
	c := &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire,
			IssuedAt:  issue,
			Issuer:    "0x" + key.PublicKey().HexString(),
			Subject:   subject,
		},
		Scope: scope,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, c)
	return token.SignedString(key.EcdsaPrivateKey())
}

func Test_simpleToken(t *testing.T) {
	require := require.New(t)
	key, _ := crypto.HexStringToPrivateKey("bc145bb9f00d55a3571e22660ef5fd1bfa596e272b80add2919735b82c273004")
	issue := time.Now().Unix()
	expire := time.Now().Add(time.Hour * 240).Unix()
	subject := "http://example.come/1234"
	scopes := []string{"write:pea", "read:pea", "delete:pea", "create:pods", "read:pods", "delete:pods"}
	for _, scope := range scopes {
		token, err := SignJWT(issue, expire, subject, scope, key)
		require.NoError(err)
		t.Logf("scope: %s,token : %s", scope, token)
	}
}

func Test_Sign(t *testing.T) {
	require := require.New(t)
	priKey, err := crypto.HexStringToPrivateKey("bc145bb9f00d55a3571e22660ef5fd1bfa596e272b80add2919735b82c273004")

	require.NoError(err)
	claim := &Claims{}

	jwtStrings := []string{
		"eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZSI6ImRlbGV0ZTpwb2RzIiwiZXhwIjoxNjA3NjY5MjQ5LCJpYXQiOjE2MDY4MDUyNDksImlzcyI6IjB4MDRmY2VkMWJhYTQxMDg2ZmU1NjhhNDhhMzM4YWYxMDRlZTE5Nzc4MDQzZDk4YzIyNjU1NzM0ZGM4ODU5MDg2MWIyNjlkZTE4NzNiN2Y4ZmFjNGRhODY3YjI0YTdjNzQ3MzlmZjNkNDZmNmQwM2M5ZGFiOGM3MTA2YmFmYjk3YTgwOSIsInN1YiI6Imh0dHA6Ly9leGFtcGxlLmNvbWUvMTIzNCJ9.nNjSIpO0OykwPXZLRFuKhMXg0tbxsCGjF3wMEiuZr2lK9-ynM3Ct4-CTxnnM5mI4HETddKI9Mnh74PpkJJibjA",
		"eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZSI6IndyaXRlOnBlYSIsImV4cCI6MTYwNzY2OTI0OSwiaWF0IjoxNjA2ODA1MjQ5LCJpc3MiOiIweDA0ZmNlZDFiYWE0MTA4NmZlNTY4YTQ4YTMzOGFmMTA0ZWUxOTc3ODA0M2Q5OGMyMjY1NTczNGRjODg1OTA4NjFiMjY5ZGUxODczYjdmOGZhYzRkYTg2N2IyNGE3Yzc0NzM5ZmYzZDQ2ZjZkMDNjOWRhYjhjNzEwNmJhZmI5N2E4MDkiLCJzdWIiOiJodHRwOi8vZXhhbXBsZS5jb21lLzEyMzQifQ.PmS28Z2lKcNaZWsaZQJnO-Po1JRVj5oG28fQ1sT9wp6QoGsrSDpVChZHUIveG6V8sfTi_lqOcgcTgvkTqinFDg",
	}
	for _, jwtString := range jwtStrings {
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
		require.NoError(err)
		require.Equal(true, token.Valid)
		str, err := token.SigningString()
		require.NoError(err)
		sig, err := token.Method.Sign(str, priKey.EcdsaPrivateKey())
		require.NoError(err)
		require.NotEqual(sig, token.Signature)

	}
}
