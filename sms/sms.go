package sms

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"strings"
)

type Sms struct {
	client *dysmsapi.Client    // 短信推送客户端
	valid  *validator.Validate // 参数验证器
}

type Params struct {
	Phones        string `validate:"required=true"` // 被推送短信的手机号，格式：15612345678,1341234567
	TemplateID    string `validate:"required=true"` // 短信模板ID
	Sign          string `validate:"required=true"` // 签名
	TemplateParam string // 模板参数
}

func (as *Sms) Send(opts *Params) error {
	// 通用验证
	if err := as.valid.Struct(opts); err != nil {
		return errors.Wrap(err, "通用参数验证失败")
	}

	// 验证签名，阿里云短信推送必须传入签名
	if err := as.valid.Var(opts.Sign, "required=true"); err != nil {
		return errors.Wrap(err, "签名验证失败")
	}

	// 根据手机号的多少判断是单发，但是群发
	phones := strings.Split(opts.Phones, ",")
	if len(phones) > 1 {
		return as.SendBatchSms(opts) // 群发
	} else {
		return as.SendSms(opts) // 单发
	}
}

// 单个手机号发送短信
func (as *Sms) SendSms(opts *Params) error {
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	request.SignName = opts.Sign
	request.PhoneNumbers = opts.Phones
	request.TemplateCode = opts.TemplateID
	request.TemplateParam = opts.TemplateParam

	response, err := as.client.SendSms(request)
	if err != nil {
		return errors.Wrap(err, "短信单发推送失败")
	}

	fmt.Println("短信单发结果: ")
	fmt.Printf("BizId - %#v\n", response.BizId)
	fmt.Printf("Code - %#v\n", response.Code)
	fmt.Printf("Message - %#v\n", response.Message)
	fmt.Printf("RequestId - %#v\n", response.RequestId)

	return nil
}

func (as *Sms) SendBatchSms(opts *Params) error {
	request := dysmsapi.CreateSendBatchSmsRequest()
	request.Scheme = "https"

	var phoneBuilder strings.Builder
	phoneBuilder.WriteString("[")
	phones := strings.Split(opts.Phones, ",")
	for _, v := range phones {
		phoneBuilder.WriteString(v)
	}
	phoneBuilder.WriteString("]")

	var signBuilder strings.Builder
	signBuilder.WriteString("[")
	signs := strings.Split(opts.Sign, ",")
	for _, v := range signs {
		signBuilder.WriteString(v)
	}
	signBuilder.WriteString("]")

	request.PhoneNumberJson = phoneBuilder.String()
	request.SignNameJson = signBuilder.String()

	request.TemplateCode = opts.TemplateID
	request.TemplateParamJson = opts.TemplateParam // ep："[{\"code\":\"666666\"}]" 是一个json的对象数组

	response, err := as.client.SendBatchSms(request)
	if err != nil {
		return errors.Wrap(err, "短信群发推送失败")
	}

	fmt.Println("短信群发结果: ")
	fmt.Printf("BizId - %#v\n", response.BizId)
	// 根据code是否为ok判断是否发送成功
	// 其他值可以在地址 https://help.aliyun.com/document_detail/101346.html?spm=a2c4g.11186623.2.14.8042128eDUcH7U 找到原因。
	fmt.Printf("Code - %#v\n", response.Code)
	fmt.Printf("Message - %#v\n", response.Message)
	fmt.Printf("RequestId - %#v\n", response.RequestId)

	return nil
}

// 创建验证对象
func NewValidate() (*validator.Validate, error) {
	valid := validator.New()
	// 注册自定义验证规则
	return valid, nil
}

// 注意：这里在创建aliSms的时候，采用就近赋值原则，即：创建一个字段，赋值一个字段。
func NewSms(AccessKeyId, AccessSecret, RegionId string) (*Sms, error) {
	var ali Sms

	// 验证器
	valid, err := NewValidate()
	if err != nil {
		return nil, err
	}
	ali.valid = valid

	// 短信推送客户端
	client, err := dysmsapi.NewClientWithAccessKey(RegionId, AccessKeyId, AccessSecret)
	if err != nil {
		return nil, err
	}
	ali.client = client

	return &ali, nil
}
