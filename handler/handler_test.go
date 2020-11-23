package handler

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/iotexproject/phoenix-gem/auth"
	"github.com/stretchr/testify/require"
)

func simpleToken() (string, error) {
	mySigningKey := []byte("123456789")

	// Create the Claims
	claims := auth.Claims{
		Scope: "bar",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

func Test_simpleToken(t *testing.T) {
	require := require.New(t)
	token, err := simpleToken()
	require.NoError(err)
	// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0.QUCEXUBDPNaLC_soGYb0g9ErI3VDPYLpIzGor0WJgIo
	t.Logf("%v", token)
}
