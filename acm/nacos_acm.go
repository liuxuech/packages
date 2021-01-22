package acm

type nacosACM struct {
	options Options
}

func (a nacosACM) GetConfig() (string, error) {
	panic("implement me")
}

func newRegistry(opts ...Option) ACM {
	var nacos nacosACM

	// 设置默认值
	nacos.options.TimeoutMs = 5 * 1000
	nacos.options.IpAddr = "127.0.0.1"
	nacos.options.Port = 8848

	for _, o := range opts {
		o(&nacos.options)
	}

	// namespaceId 没有默认值，但是又是必须的值，所以需要验证

	return &nacosACM{}
}
