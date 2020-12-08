package store

type Options struct {
	Database string // ep: goslam.db  Database(goslam)
	Dir      string
}

type ReadOptions struct {
}
type WriteOptions struct {
}
type DeleteOptions struct {
}
type ListOptions struct {
}

type Option func(*Options)
type ReadOption func(*ReadOptions)
type WriteOption func(*WriteOptions)
type DeleteOption func(*DeleteOptions)
type ListOption func(*ListOptions)
