package ali

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/go-playground/validator/v10"
	"github.com/liuxuech/packages/sms"
	"github.com/pkg/errors"
)

type aliSms struct {
	opts   *Options            // 配置
	client *dysmsapi.Client    // 短信推送客户端
	valid  *validator.Validate // 参数验证器
}

func (as *aliSms) Send(opts *sms.MessageOption) error {
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	// 验证签名，阿里云短信推送必须传入签名
	if err := as.valid.Var(opts.Sign, "required=true"); err != nil {
		return errors.Wrap(err, "签名验证失败")
	}

	request.SignName = opts.Sign
	request.PhoneNumbers = opts.Phones
	request.TemplateCode = opts.TemplateID
	request.TemplateParam = opts.TemplateParam

	response, err := as.client.SendSms(request)
	if err != nil {
		return errors.Wrap(err, "短信推送失败")
	}

	fmt.Println("短信发送结果: ")
	fmt.Printf("BizId - %#v\n", response.BizId)
	fmt.Printf("Code - %#v\n", response.Code)
	fmt.Printf("Message - %#v\n", response.Message)
	fmt.Printf("RequestId - %#v\n", response.RequestId)

	return nil
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

	as.opts = &options

	if err := as.init(&options); err != nil {
		return nil, err
	}

	return &as, nil
}
