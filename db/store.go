package db

import (
	"context"
)

type (
	storeCtxKey struct{}

	Store interface {
		// Name is the name of the store
		Name() string

		// Region is the region where the storage is located
		Region() string

		// Endpoint is the endpoint URL
		Endpoint() string

		// AccessKey is the key/ID to access the endpoint
		AccessKey() string

		// AccessToken is the token to access the endpoint
		AccessToken() string
	}

	store struct {
		name     string
		region   string
		endpoint string
		key      string
		token    string
	}
)

// WithStoreCtx add StoreCtx into context.
func WithStoreCtx(ctx context.Context, s Store) context.Context {
	return context.WithValue(ctx, storeCtxKey{}, s)
}

// GetBlockCtx gets BlockCtx
func GetStoreCtx(ctx context.Context) (Store, bool) {
	store, ok := ctx.Value(storeCtxKey{}).(Store)
	return store, ok
}

func (s *store) Name() string {
	return s.name
}

func (s *store) Region() string {
	return s.region
}

func (s *store) Endpoint() string {
	return s.endpoint
}

func (s *store) AccessKey() string {
	return s.key
}

func (s *store) AccessToken() string {
	return s.token
}
