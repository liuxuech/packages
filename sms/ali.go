package sms

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/go-playground/validator/v10"
	"github.com/liuxuech/packages/base"
)

type aliSms struct {
	opts   *Options
	client *dysmsapi.Client
	valid  *validator.Validate
}

func (as *aliSms) Send(opts *SingleMsg) error {
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	request.PhoneNumbers = opts.TargetPhone
	request.SignName = opts.Sign
	request.TemplateCode = opts.TemplateCode
	request.TemplateParam = opts.TemplateParam

	response, err := as.client.SendSms(request)
	if err != nil {
		return err
	}

	fmt.Println("短信发送结果: ")
	fmt.Printf("BizId - %#v\n", response.BizId)
	fmt.Printf("Code - %#v\n", response.Code)
	fmt.Printf("Message - %#v\n", response.Message)
	fmt.Printf("RequestId - %#v\n", response.RequestId)

	return nil
}

func (as *aliSms) SendBatch() {
	panic("implement me")
}

func (as *aliSms) validate(params interface{}) error {
	return as.valid.Struct(params)
}

func (as *aliSms) PhoneCheck(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	return base.PhoneCheck(phone)
}

func (as *aliSms) init(opts *Options) error {
	// 验证器
	as.valid = validator.New()
	// 注册自定义验证规则
	if err := as.valid.RegisterValidation("PhoneCheck", as.PhoneCheck); err != nil {
		return err
	}

	if err := as.validate(opts); err != nil {
		return err
	}

	// 短信推送客户端
	smsClient, err := dysmsapi.NewClientWithAccessKey(as.opts.RegionId, as.opts.AccessKeyId, as.opts.AccessSecret)
	if err != nil {
		return err
	}

	as.client = smsClient

	return nil
}

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

func NewSms(opts ...Option) (Sms, error) {
	var (
		options Options
		sms     aliSms
	)
	for _, o := range opts {
		o(&options)
	}
	if err := sms.init(&options); err != nil {
		return nil, err
	}

	return &sms, nil
}
