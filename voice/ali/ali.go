package ali

import (
	"fmt"
	"github.com/pkg/errors"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dyvmsapi"
	"github.com/go-playground/validator/v10"
	"github.com/liuxuech/packages/voice"
)

type aliVoice struct {
	opts   *Options
	client *dyvmsapi.Client
	valid  *validator.Validate
}

func (av *aliVoice) Call(opts *voice.Options) error {
	// 通用验证
	if err := av.valid.Struct(opts); err != nil {
		return errors.Wrap(err, "通用参数验证失败")
	}
	return av.SingleCallByTts(opts)
}

func (av *aliVoice) SingleCallByTts(opts *voice.Options) error {
	client, err := dyvmsapi.NewClientWithAccessKey(av.opts.RegionId, av.opts.AccessKeyId, av.opts.AccessSecret)
	if err != nil {
		return err
	}

	request := dyvmsapi.CreateSingleCallByTtsRequest()
	request.Scheme = "https"

	request.TtsCode = opts.TtsCode
	request.CalledNumber = opts.CalledNumber

	request.CalledShowNumber = opts.CalledShowNumber
	request.TtsParam = opts.TtsParam

	response, err := client.SingleCallByTts(request)
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

func NewVoice(opts ...Option) (voice.Voice, error) {
	var ali aliVoice

	var options Options
	for _, o := range opts {
		o(&options)
	}
	ali.opts = &options

	// 验证器
	valid, err := NewValidate()
	if err != nil {
		return nil, err
	}
	if err := valid.Struct(options); err != nil {
		return nil, err
	}
	ali.valid = valid

	// 短信推送客户端
	client, err := dyvmsapi.NewClientWithAccessKey(ali.opts.RegionId, ali.opts.AccessKeyId, ali.opts.AccessSecret)
	if err != nil {
		return nil, err
	}
	ali.client = client

	return &ali, nil
}
