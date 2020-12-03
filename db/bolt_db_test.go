package db

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestBoltDB(t *testing.T) {
	r := require.New(t)

	testFile, err := ioutil.TempFile(os.TempDir(), "test-bolt")
	path := testFile.Name()
	r.NoError(err)
	testFile.Close()
	defer func() {
		r.NoError(os.Remove(path))
	}()

	db := NewBoltDB(path)
	r.NotNil(db)
	ctx := context.Background()
	r.NoError(db.Start(ctx))
	defer func() {
		r.NoError(db.Stop(ctx))
	}()

	dbTests := []struct {
		ns   string
		k, v []byte
		err  error
	}{
		{"5NJ2Hqv",
			[]byte("JtQTAme2SKJzXVs"), []byte("nCb6ZLdB7NRaqsm91JtQTAme2SKJzXVsdPIPkyJr1MU"),
			ErrBucketNotExist,
		},
		{"5NJ2Hqv",
			[]byte("exUALhrxXi3DcLg2"), []byte("nCb6ZLdB7NRaqsm91JtQTAme2SKJzXVsdPIPkyJr1MU"),
			ErrNotExist,
		},
		{"5NJ2Hqv",
			[]byte("JtQTAme2SKJzXVs"), []byte("exUALhrxXi3DcLg2Ts+ymUAY5y4NvMg=tzjeg+3xrA5"),
			nil,
		},
		{"skON4iLHjc5H38",
			[]byte("JtQTAme2SKJzXVs"), []byte("exUALhrxXi3DcLg2Ts+ymUAY5y4NvMg=tzjeg+3xrA5"),
			ErrBucketNotExist,
		},
		{"skON4iLHjc5H38",
			[]byte("exUALhrxXi3DcLg2"), []byte{},
			ErrNotExist,
		},
	}

	for _, e := range dbTests {
		v, err := db.Get(e.ns, e.k)
		if e.err != nil {
			r.Equal(e.err, errors.Cause(err))
			r.Nil(v)
		}

		r.NoError(db.Put(e.ns, e.k, e.v))
		v, err = db.Get(e.ns, e.k)
		r.NoError(err)
		r.Equal(e.v, v)

		r.NoError(db.Delete(e.ns, e.k))
		v, err = db.Get(e.ns, e.k)
		r.Equal(ErrNotExist, errors.Cause(err))
	}
}
