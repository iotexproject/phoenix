package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetStoreCtx(t *testing.T) {
	r := require.New(t)

	s := &store{
		name:     "s3",
		endpoint: "s3/abcdef",
		token:    "E3ru+11k8xSBh+hMPgOLZmtrrCbhqsmaPHjLKYnJCaQ=",
	}

	ctx := WithStoreCtx(context.Background(), s)
	r.NotNil(ctx)
	s1, ok := GetStoreCtx(ctx)
	r.True(ok)
	r.Equal(s.name, s1.Name())
	r.Equal(s.endpoint, s1.Endpoint())
	r.Equal(s.token, s1.AccessToken())
}
