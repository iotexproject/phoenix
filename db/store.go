package db

import (
	"context"
)

type (
	storeCtxKey struct{}

	Store interface {
		// Name is the name of the store
		Name() string

		// Endpoint is the endpoint URL
		Endpoint() string

		// AccessToken is the token to access the endpoint
		AccessToken() string
	}
	
	store struct {
		name     string
		endpoint string
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

func (s *store) Endpoint() string {
	return s.endpoint
}

func (s *store) AccessToken() string {
	return s.token
}
