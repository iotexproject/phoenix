package midware

import (
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/iotexproject/go-pkgs/crypto"
	"github.com/iotexproject/iotex-antenna-go/v2/jwt"
	"github.com/iotexproject/phoenix-gem/auth"
	"github.com/iotexproject/phoenix-gem/db"
)

func TestCredential(t *testing.T) {
	r := require.New(t)

	testFile, err := ioutil.TempFile(os.TempDir(), "test-bolt")
	path := testFile.Name()
	r.NoError(err)
	testFile.Close()
	defer func() {
		r.NoError(os.Remove(path))
	}()

	db := db.NewBoltDB(path)
	r.NotNil(db)
	ctx := context.Background()
	r.NoError(db.Start(ctx))
	defer func() {
		r.NoError(db.Stop(ctx))
	}()

	c := NewCredential(db)
	user := "71099b90dDC322a773115295f560bD1Af02f777d"
	tag := "s3"
	s := auth.NewStore(
		"s3",
		"us-east",
		"s3/abcdef",
		"AKIAIOSFODNN7EXAMPLE",
		"E3ru+11k8xSBh+hMPgOLZmtrrCbhqsmaPHjLKYnJCaQ=",
	)
	r.NoError(c.PutStore(user, tag, s))
	s1, err := c.GetStore(user, tag)
	r.NoError(err)
	r.Equal(s, s1)
}

func Test_CredentialPutTestData(t *testing.T) {
	r := require.New(t)
	path := os.Getenv("dbpath")
	if path == "" {
		r.Empty(path)
		return
	}
	db := db.NewBoltDB(path)
	r.NotNil(db)
	ctx := context.Background()
	r.NoError(db.Start(ctx))

	c := NewCredential(db)
	user := "6a26b3056679e6adf079350c778ae7ab71a287fa"
	tag := "s3"
	s := auth.NewStore(
		"s3",
		"us-east",
		"http://localhost:9001",
		"AKIAIOSFODNN7EXAMPLE",
		"wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
	)
	r.NoError(c.PutStore(user, tag, s))
	s1, err := c.GetStore(user, tag)
	r.NoError(err)
	r.Equal(s, s1)
	r.NoError(db.Stop(ctx))
}

func Test_JwtData(t *testing.T) {
	r := require.New(t)
	key, _ := crypto.HexStringToPrivateKey("bc145bb9f00d55a3571e22660ef5fd1bfa596e272b80add2919735b82c273004")
	issue := time.Now().Unix()
	expire := time.Now().Add(time.Hour * 240).Unix()
	subject := "s3"
	scopes := []string{jwt.CREATE, jwt.DELETE, jwt.UPDATE, jwt.READ}
	for _, scope := range scopes {
		token, err := jwt.SignJWT(issue, expire, subject, scope, key)
		r.NoError(err)
		t.Logf("scope: %s, Issuer: %s,token : %s", scope, key.PublicKey().HexString(), token)
	}
}
