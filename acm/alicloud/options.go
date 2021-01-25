package alicloud

import (
	"github.com/liuxuech/packages/acm"
)

const (
	accessKey = "AccessKey"
	secretKey = "SecretKey"
)

func WithAccessKey(access string) acm.Option {
	return func(options *acm.Options) {
		options.Others[accessKey] = access
	}
}

func WithSecretKey(secret string) acm.Option {
	return func(options *acm.Options) {
		options.Others[secretKey] = secret
	}
}
