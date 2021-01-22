package acm

var DefaultACM = NewACM()

type ACM interface {
	GetConfig() (string, error)
}

func NewACM(opts ...Option) ACM {
	return newRegistry(opts...)
}

func GetConfig() (string, error) {
	return DefaultACM.GetConfig()
}
