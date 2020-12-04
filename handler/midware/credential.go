package midware

import (
	"net/http"

	"github.com/iotexproject/go-pkgs/crypto"
	"github.com/pkg/errors"

	"github.com/iotexproject/phoenix-gem/auth"
	"github.com/iotexproject/phoenix-gem/db"
)

type (
	Credential interface {
		// GetStore returns user's store according to tag
		GetStore(string, string) (auth.Store, error)

		// PutStore puts user's store into db
		PutStore(string, string, auth.Store) error

		// DoCredential
		DoCredential(http.Handler) http.Handler
	}

	credential struct {
		db.KVStore
	}
)

func NewCredential(kv db.KVStore) Credential {
	return &credential{
		KVStore: kv,
	}
}

func (c *credential) GetStore(user, tag string) (auth.Store, error) {
	bytes, err := c.Get(user, []byte(tag))
	if err != nil {
		return nil, err
	}
	return auth.DeserializeToStore(bytes)
}

func (c *credential) PutStore(user, tag string, store auth.Store) error {
	bytes, err := store.Serialize()
	if err != nil {
		return err
	}
	return c.Put(user, []byte(tag), bytes)
}

func (c *credential) DoCredential(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := ctx.Value(auth.TokenCtxKey).(*auth.Claims)
		if !ok {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// trustor is the address that registers endpoint with us
		trustor, err := crypto.HexStringToPublicKey(claims.Issuer)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// check trustor's storage endpoint
		name := trustor.Address().Hex()[2:] // remove 0x prefix
		store, err := c.GetStore(name, claims.Subject)
		switch errors.Cause(err) {
		case nil:
			// pass the store to next handler
			ctx := auth.WithStoreCtx(ctx, store)
			next.ServeHTTP(w, r.WithContext(ctx))
		case db.ErrBucketNotExist, db.ErrNotExist:
			http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	})
}
