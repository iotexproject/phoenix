package auth

import (
	"github.com/iotexproject/phoenix-gem/db"
)

type (
	Credential interface {
		// GetStore returns user's store according to tag
		GetStore(string, string) (Store, error)

		// PutStore puts user's store into db
		PutStore(string, string, Store) error

		// DelStore deletes user's store from db
		DelStore(string, string) error
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

func (c *credential) GetStore(user, tag string) (Store, error) {
	bytes, err := c.Get(user, []byte(tag))
	if err != nil {
		return nil, err
	}
	return DeserializeToStore(bytes)
}

func (c *credential) PutStore(user, tag string, store Store) error {
	bytes, err := store.Serialize()
	if err != nil {
		return err
	}
	return c.Put(user, []byte(tag), bytes)
}

func (c *credential) DelStore(user, tag string) error {
	return c.Delete(user, []byte(tag))
}
