package ali

import (
	"github.com/go-playground/validator/v10"
	"github.com/liuxuech/packages/base"
)

type Options struct {
	RegionId     string `validate:"required=true"`
	AccessKeyId  string `validate:"required=true"`
	AccessSecret string `validate:"required=true"`
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

func NewValidate() (*validator.Validate, error) {
	valid := validator.New()
	// 注册自定义验证规则
	if err := valid.RegisterValidation("PhoneCheck", PhoneCheck); err != nil {
		return nil, err
	}
	return valid, nil
}

func PhoneCheck(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	return base.PhoneCheck(phone)
}
