package auth

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/iotexproject/phoenix/db"
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

	d := db.NewBoltDB(path)
	r.NotNil(d)
	ctx := context.Background()
	r.NoError(d.Start(ctx))
	defer func() {
		r.NoError(d.Stop(ctx))
	}()

	c := NewCredential(d)
	user := "71099b90dDC322a773115295f560bD1Af02f777d"
	tag := "s3"
	s := NewStore(
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

	r.NoError(c.DelStore(user, tag))
	s1, err = c.GetStore(user, tag)
	r.Equal(db.ErrNotExist, errors.Cause(err))
	r.Nil(s1)
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
	s := NewStore(
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
