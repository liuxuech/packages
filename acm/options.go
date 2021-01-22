package acm

// 公共的配置参数
type Options struct {
	IpAddr      string
	Port        uint64
	TimeoutMs   uint64 // 超时时间，单位 毫秒
	NamespaceId string // 命名空间，ep：dev、test、prod
}

type Option func(*Options)

func Host(ipAddr string) Option {
	return func(options *Options) {
		options.IpAddr = ipAddr
	}
}

func Port(post uint64) Option {
	return func(options *Options) {
		options.Port = post
	}
}

func TimeoutMs(time uint64) Option {
	return func(options *Options) {
		options.TimeoutMs = time
	}
}

func NamespaceId(nsId string) Option {
	return func(options *Options) {
		options.NamespaceId = nsId
	}
}
