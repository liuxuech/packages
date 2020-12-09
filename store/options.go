package store

import "context"

type Options struct {
	Database string
	Table    string
	Context  context.Context
}

func WithDatabase(database string) Option {
	return func(opts *Options) {
		opts.Database = database
	}
}

func WithTable(table string) Option {
	return func(opts *Options) {
		opts.Table = table
	}
}

type ReadOptions struct {
	Table string
}

func ReadFrom(table string) ReadOption {
	return func(opts *ReadOptions) {
		opts.Table = table
	}
}

type WriteOptions struct {
	Table string
}

func WriteTo(table string) WriteOption {
	return func(opts *WriteOptions) {
		opts.Table = table
	}
}

type DeleteOptions struct {
	Table string
}

func DeleteFrom(table string) DeleteOption {
	return func(opts *DeleteOptions) {
		opts.Table = table
	}
}

type ListOptions struct {
}

type Option func(*Options)
type ReadOption func(*ReadOptions)
type WriteOption func(*WriteOptions)
type DeleteOption func(*DeleteOptions)
type ListOption func(*ListOptions)
