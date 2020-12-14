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

	// 通用验证
	if err := as.valid.Struct(opts); err != nil {
		return errors.Wrap(err, "通用参数验证失败")
	}

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

func (as *aliSms) init(opts *Options) (err error) {

	// 短信推送客户端
	smsClient, err := dysmsapi.NewClientWithAccessKey(as.opts.RegionId, as.opts.AccessKeyId, as.opts.AccessSecret)
	if err != nil {
		return err
	}

	as.client = smsClient

	return nil
}

// 创建验证对象
func NewValidate() (*validator.Validate, error) {
	valid := validator.New()
	// 注册自定义验证规则
	return valid, nil
}

func NewSms(opts ...Option) (sms.Sms, error) {
	var (
		ali     aliSms
		options Options
	)

	for _, o := range opts {
		o(&options)
	}

	// 验证器
	valid, err := NewValidate()
	if err != nil {
		return nil, err
	}
	if err := valid.Struct(options); err != nil {
		return nil, err
	}

	// 短信推送客户端
	client, err := dysmsapi.NewClientWithAccessKey(options.RegionId, options.AccessKeyId, options.AccessSecret)
	if err != nil {
		return nil, err
	}

	ali.opts = &options
	ali.valid = valid
	ali.client = client

	return &ali, nil
}
