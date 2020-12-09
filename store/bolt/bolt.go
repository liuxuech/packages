package bolt

import (
	"github.com/liuxuech/packages/store"
	bolt "go.etcd.io/bbolt"
)

type boltStore struct {
	db      *bolt.DB
	options store.Options
}

func (b *boltStore) Read(key string, opts ...store.ReadOption) (*store.Record, error) {
	panic("implement me")
}

func (b *boltStore) Write(r *store.Record, opts ...store.WriteOption) error {
	panic("implement me")
}

func (b *boltStore) Delete(r *store.Record, opts ...store.DeleteOption) error {
	panic("implement me")
}

func (b *boltStore) List(opts ...store.ListOption) ([]string, error) {
	panic("implement me")
}

func (b *boltStore) Close() error {
	panic("implement me")
}

func NewStore() store.Store {
	return &boltStore{}
}
