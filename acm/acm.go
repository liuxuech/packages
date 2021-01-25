package acm

type ACM interface {
	Init() error
	GetConfig() (string, error)
}

func NewACM(opts ...Option) ACM {
	return newACM(opts...)
}
