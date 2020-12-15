package store

type Options struct {
	DBPath string // 数据文件存放的位置
	DBFile string // 数据文件
	Bucket string // 默认值的bucket
}

func WithDBPath(dbPath string) Option {
	return func(opts *Options) {
		opts.DBPath = dbPath
	}
}

func WithDBFile(dbFile string) Option {
	return func(opts *Options) {
		opts.DBFile = dbFile
	}
}

func WithBucket(bucket string) Option {
	return func(opts *Options) {
		opts.Bucket = bucket
	}
}

type ReadOptions struct {
	Bucket string
}

func ReadFrom(bucket string) ReadOption {
	return func(opts *ReadOptions) {
		opts.Bucket = bucket
	}
}

type WriteOptions struct {
	Bucket string
}

func WriteTo(bucket string) WriteOption {
	return func(opts *WriteOptions) {
		opts.Bucket = bucket
	}
}

type DeleteOptions struct {
	Bucket string
}

func DeleteFrom(bucket string) DeleteOption {
	return func(opts *DeleteOptions) {
		opts.Bucket = bucket
	}
}

type ListOptions struct {
}

type Option func(*Options)
type ReadOption func(*ReadOptions)
type WriteOption func(*WriteOptions)
type DeleteOption func(*DeleteOptions)
type ListOption func(*ListOptions)
