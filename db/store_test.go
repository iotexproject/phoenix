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
		region:   "us-east",
		endpoint: "s3/abcdef",
		key:      "AKIAIOSFODNN7EXAMPLE",
		token:    "E3ru+11k8xSBh+hMPgOLZmtrrCbhqsmaPHjLKYnJCaQ=",
	}

	ctx := WithStoreCtx(context.Background(), s)
	r.NotNil(ctx)
	s1, ok := GetStoreCtx(ctx)
	r.True(ok)
	r.Equal(s.name, s1.Name())
	r.Equal(s.region, s1.Region())
	r.Equal(s.endpoint, s1.Endpoint())
	r.Equal(s.key, s1.AccessKey())
	r.Equal(s.token, s1.AccessToken())
}
