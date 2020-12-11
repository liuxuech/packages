package ali

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/go-playground/validator/v10"
	"github.com/liuxuech/packages/sms"
)

type aliSms struct {
	opts   *Options
	client *dysmsapi.Client
	valid  *validator.Validate
}

func (as *aliSms) Send(opts *sms.SingleMsg) error {
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

func (as *aliSms) init(opts *Options) error {
	// 验证器
	valid, err := NewValidate()
	if err != nil {
		return nil
	}
	as.valid = valid

	if err := as.valid.Struct(opts); err != nil {
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

func NewSms(opts ...Option) (sms.Sms, error) {
	var (
		options Options
		as      aliSms
	)
	for _, o := range opts {
		o(&options)
	}
	if err := as.init(&options); err != nil {
		return nil, err
	}

	return &as, nil
}
