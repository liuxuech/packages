package ali

type Options struct {
	AccessKeyId  string `validate:"required=true"`
	AccessSecret string `validate:"required=true"`
	RegionId     string `validate:"required=true"`
}

type Option func(*Options)

func WithRegionId(regionId string) Option {
	return func(options *Options) {
		options.RegionId = regionId
	}
}

func WithAccessKeyId(accessKeyId string) Option {
	return func(options *Options) {
		options.AccessKeyId = accessKeyId
	}
}

func WithAccessSecret(accessSecret string) Option {
	return func(options *Options) {
		options.AccessSecret = accessSecret
	}
}
