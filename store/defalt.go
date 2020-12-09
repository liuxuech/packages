package store

import (
	"context"
	"encoding/json"
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

func (bs *boltStore) Read(key string, opts ...ReadOption) (*Record, error) {
	readOpts := ReadOptions{Table: DefaultTable}
	for _, o := range opts {
		o(&readOpts)
	}

	var (
		value  []byte
		record Record
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

func (bs *boltStore) Write(r *Record, opts ...WriteOption) error {
	writeOpts := WriteOptions{Table: DefaultTable}
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

func (bs *boltStore) Delete(key string, opts ...DeleteOption) error {
	deleteOpts := DeleteOptions{Table: DefaultTable}
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

func (bs boltStore) List(opts ...ListOption) ([]string, error) {
	return nil, nil
}

func (bs *boltStore) Close() error {
	return bs.db.Close()
}

// 设置数据文件存放的位置
func WithDir(dir string) Option {
	return func(opts *Options) {
		opts.Context = context.WithValue(opts.Context, dirKey{}, dir)
	}
}

func NewStore(opts ...Option) (Store, error) {
	options := Options{
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
