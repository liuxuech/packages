package store

import (
	"errors"
)

var (
	ErrNotFound = errors.New("not found")
)

type Store interface {
	Read(key string, opts ...ReadOption) (*Record, error)
	Write(r *Record, opts ...WriteOption) error
	Delete(key string, opts ...DeleteOption) error
	List(opts ...ListOption) ([]string, error) // 返回 keys
	Close() error
}

type Record struct {
	Key   string `json:"key"`
	Value []byte `json:"value"`
}

var defaultStore Store

func Read(key string, opts ...ReadOption) (*Record, error) {
	return defaultStore.Read(key, opts...)
}

func Write(r *Record, opts ...WriteOption) error {
	return defaultStore.Write(r, opts...)
}

func Delete(key string, opts ...DeleteOption) error {
	return defaultStore.Delete(key, opts...)
}

func Close() error {
	return defaultStore.Close()
}

func InitStore(dbPath, dbFile string) (err error) {
	if defaultStore, err = NewStore(WithDBPath(dbPath), WithDBFile(dbFile)); err != nil {
		return err
	}
	return nil
}
