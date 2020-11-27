package auth

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/iotexproject/go-pkgs/crypto"
	"github.com/stretchr/testify/require"
)

// SignJWT creates a JWT
func SignJWT(issue, expire int64, subject string, key crypto.PrivateKey) (string, error) {
	claim := &jwt.StandardClaims{
		ExpiresAt: expire,
		IssuedAt:  issue,
		Issuer:    "0x" + key.PublicKey().HexString(),
		Subject:   subject,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claim)
	return token.SignedString(key.EcdsaPrivateKey())
}

func simpleToken() (string, error) {
	key, _ := crypto.HexStringToPrivateKey("bc145bb9f00d55a3571e22660ef5fd1bfa596e272b80add2919735b82c273004")
	issue := time.Now().Unix()
	expire := time.Now().Add(time.Hour * 24).Unix()
	subject := "read:pods"
	return SignJWT(issue, expire, subject, key)
}

func Test_simpleToken(t *testing.T) {
	require := require.New(t)
	token, err := simpleToken()
	require.NoError(err)
	// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0.QUCEXUBDPNaLC_soGYb0g9ErI3VDPYLpIzGor0WJgIo
	t.Logf("%v", token)
}
