package auth

import (
	"context"

	"github.com/golang/protobuf/proto"

	"github.com/iotexproject/phoenix-gem/auth/storepb"
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

		// Serialize returns the serialized byte-stream
		Serialize() ([]byte, error)
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

func NewStore(name, region, endpoint, key, token string) Store {
	return &store{
		name:     name,
		region:   region,
		endpoint: endpoint,
		key:      key,
		token:    token,
	}
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

func (s *store) Serialize() ([]byte, error) {
	pb := &storepb.Store{
		Name:     s.name,
		Region:   s.region,
		Endpoint: s.endpoint,
		Key:      s.key,
		Token:    s.token,
	}
	return proto.Marshal(pb)
}

func DeserializeToStore(buf []byte) (*store, error) {
	pb := &storepb.Store{}
	if err := proto.Unmarshal(buf, pb); err != nil {
		return nil, err
	}
	return &store{
		name:     pb.Name,
		region:   pb.Region,
		endpoint: pb.Endpoint,
		key:      pb.Key,
		token:    pb.Token,
	}, nil
}
