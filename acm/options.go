package acm

// 公共的配置参数
type Options struct {
	Host     string
	Port     uint64
	LogDir   string
	CacheDir string

	NamespaceId string // 命名空间，ep：dev、test、prod

	// 其他实现所需的配置可以存放在这里
	Others map[string]interface{}
}

type Option func(*Options)

func WithHost(host string) Option {
	return func(options *Options) {
		options.Host = host
	}
}

func WithPort(post uint64) Option {
	return func(options *Options) {
		options.Port = post
	}
}

func WithNamespaceId(nsId string) Option {
	return func(options *Options) {
		options.NamespaceId = nsId
	}
}

func WithCacheDir(dir string) Option {
	return func(options *Options) {
		options.CacheDir = dir
	}
}

func WithLogDir(dir string) Option {
	return func(options *Options) {
		options.LogDir = dir
	}
}
