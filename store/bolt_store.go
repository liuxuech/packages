package store

import (
	"encoding/json"
	bolt "go.etcd.io/bbolt"
	"os"
	"path/filepath"
)

const (
	DefaultDBPath = "."
	DefaultDBFile = "bolt.db"
	DefaultBucket = "bolt"
)

type boltStore struct {
	opts *Options
	db   *bolt.DB
}

func (bs *boltStore) Read(key string, opts ...ReadOption) (*Record, error) {
	readOpts := ReadOptions{Bucket: bs.opts.Bucket}
	for _, o := range opts {
		o(&readOpts)
	}

	var (
		value  []byte
		record Record
	)

	err := bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(readOpts.Bucket))
		if b == nil {
			return nil
		}
		value = b.Get([]byte(key))
		return nil
	})

	if value == nil {
		return nil, ErrNotFound
	}

	if err = json.Unmarshal(value, &record); err != nil {
		return nil, err
	}

	return &record, err
}

func (bs *boltStore) Write(r *Record, opts ...WriteOption) error {
	writeOpts := WriteOptions{Bucket: bs.opts.Bucket}
	for _, o := range opts {
		o(&writeOpts)
	}

	data, err := json.Marshal(r)
	if err != nil {
		return nil
	}

	return bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(writeOpts.Bucket))
		if b == nil {
			if b, err = tx.CreateBucketIfNotExists([]byte(writeOpts.Bucket)); err != nil {
				return err
			}
		}
		return b.Put([]byte(r.Key), data)
	})
}

func (bs *boltStore) Delete(key string, opts ...DeleteOption) error {
	deleteOpts := DeleteOptions{Bucket: bs.opts.Bucket}
	for _, o := range opts {
		o(&deleteOpts)
	}

	return bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(deleteOpts.Bucket))
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

func NewStore(opts ...Option) (Store, error) {
	var bs boltStore

	options := Options{
		DBPath: DefaultDBPath,
		DBFile: DefaultDBFile,
		Bucket: DefaultBucket,
	}
	for _, o := range opts {
		o(&options)
	}
	bs.opts = &options

	// 创建目录
	os.MkdirAll(options.DBPath, 0700)

	boltDB, err := bolt.Open(filepath.Join(options.DBPath, options.DBFile), 0700, nil)
	if err != nil {
		return nil, err
	}
	bs.db = boltDB

	return &bs, nil
}
