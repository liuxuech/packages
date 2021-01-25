package acm

var DefaultACM = NewACM()

type ACM interface {
	Init() error
	GetConfig() (string, error)
}

func NewACM(opts ...Option) ACM {
	return newACM(opts...)
}

func GetConfig() (string, error) {
	return DefaultACM.GetConfig()
}
