package voice

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dyvmsapi"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type Voice struct {
	client *dyvmsapi.Client
	valid  *validator.Validate
}

type Msg struct {
	TtsCode          string `validate:"required=true"`
	TtsParam         string
	CalledNumber     string `validate:"required=true"`
	CalledShowNumber string
}

func (av *Voice) Call(opts *Msg) error {
	// 通用验证
	if err := av.valid.Struct(opts); err != nil {
		return errors.Wrap(err, "通用参数验证失败")
	}
	return av.SingleCallByTts(opts)
}

func (av *Voice) SingleCallByTts(opts *Msg) error {
	request := dyvmsapi.CreateSingleCallByTtsRequest()
	request.Scheme = "https"

	request.TtsCode = opts.TtsCode
	request.CalledNumber = opts.CalledNumber

	request.CalledShowNumber = opts.CalledShowNumber
	request.TtsParam = opts.TtsParam

	response, err := av.client.SingleCallByTts(request)
	if err != nil {
		return err
	}

	fmt.Println("语音通知结果: ")
	fmt.Printf("CallId - %#v\n", response.CallId)
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

func NewVoice(AccessKeyId, AccessSecret, RegionId string) (*Voice, error) {
	var ali Voice

	// 验证器
	valid, err := NewValidate()
	if err != nil {
		return nil, err
	}
	ali.valid = valid

	// 短信推送客户端
	client, err := dyvmsapi.NewClientWithAccessKey(RegionId, AccessKeyId, AccessSecret)
	if err != nil {
		return nil, err
	}
	ali.client = client

	return &ali, nil
}
