package midware

import (
	"context"
	"github.com/iotexproject/phoenix-gem/auth"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

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
