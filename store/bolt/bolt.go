package bolt

import (
	"context"
	"encoding/json"
	"github.com/liuxuech/packages/store"
	bolt "go.etcd.io/bbolt"
	"os"
	"path/filepath"
)

const (
	DefaultPath     = "."
	DefaultDatabase = "bolt"
	DefaultTable    = "bolt" // table 对应 bucket
)

type dirKey struct{}

type boltStore struct {
	db *bolt.DB
}

func (bs *boltStore) Read(key string, opts ...store.ReadOption) (*store.Record, error) {
	readOpts := store.ReadOptions{Table: DefaultTable}
	for _, o := range opts {
		o(&readOpts)
	}

	var (
		value  []byte
		record store.Record
	)

	err := bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(readOpts.Table))
		value = b.Get([]byte(key))
		return nil
	})

	if err = json.Unmarshal(value, &record); err != nil {
		return nil, err
	}

	return &record, err
}

func (bs *boltStore) Write(r *store.Record, opts ...store.WriteOption) error {
	writeOpts := store.WriteOptions{Table: DefaultTable}
	for _, o := range opts {
		o(&writeOpts)
	}

	data, err := json.Marshal(r)
	if err != nil {
		return nil
	}

	return bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(writeOpts.Table))
		if b == nil {
			if b, err = tx.CreateBucketIfNotExists([]byte(writeOpts.Table)); err != nil {
				return err
			}
		}
		return b.Put([]byte(r.Key), data)
	})
}

func (bs *boltStore) Delete(key string, opts ...store.DeleteOption) error {
	deleteOpts := store.DeleteOptions{Table: DefaultTable}
	for _, o := range opts {
		o(&deleteOpts)
	}

	return bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(deleteOpts.Table))
		if b == nil {
			return nil
		}
		return b.Delete([]byte(key))
	})
}

func (bs boltStore) List(opts ...store.ListOption) ([]string, error) {
	return nil, nil
}

func (bs *boltStore) Close() error {
	return bs.db.Close()
}

// 设置数据文件存放的位置
func WithDir(dir string) store.Option {
	return func(opts *store.Options) {
		opts.Context = context.WithValue(opts.Context, dirKey{}, dir)
	}
}

func NewStore(opts ...store.Option) (store.Store, error) {
	options := store.Options{
		Database: DefaultDatabase,
		Table:    DefaultTable,
	}

	for _, o := range opts {
		o(&options)
	}

	var dir string
	if options.Context != nil {
		dir, _ = options.Context.Value(dirKey{}).(string)
	}
	if len(dir) == 0 {
		dir = DefaultPath
	}

	// 创建目录
	os.MkdirAll(filepath.Join(dir), 0700)

	boltDB, err := bolt.Open(filepath.Join(dir, options.Database+".db"), 0700, nil)
	if err != nil {
		return nil, err
	}

	return &boltStore{db: boltDB}, nil
}
