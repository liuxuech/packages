package store

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
