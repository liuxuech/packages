package alicloud

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/liuxuech/packages/acm"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/pkg/errors"
)

const (
	accessKey = "AccessKey"
	secretKey = "SecretKey"
)

type aliACM struct {
	options  acm.Options
	validate *validator.Validate

	client       config_client.IConfigClient
	clientConfig constant.ClientConfig
}

func (ali *aliACM) Validate() error {
	if err := ali.validate.Struct(ali.options); err != nil {
		return err
	}
	if access, ok := ali.options.Others[accessKey]; !ok || access == "" {
		return errors.New("缺少必要参数: accessKey")
	}
	if secret, ok := ali.options.Others[secretKey]; !ok || secret == "" {
		return errors.New("缺少必要参数: secretKey")
	}

	return nil
}

func (ali *aliACM) configure() error {
	access := ali.options.Others[accessKey]
	secret := ali.options.Others[secretKey]

	// access和secret 在设置的时候就指定了类型为string
	ali.clientConfig = constant.ClientConfig{
		AccessKey:   access.(string),
		SecretKey:   secret.(string),
		TimeoutMs:   5 * 1000,
		NamespaceId: ali.options.Namespace,
		Endpoint:    fmt.Sprintf("%s:%d", ali.options.Host, ali.options.Port),
	}

	client, err := clients.CreateConfigClient(map[string]interface{}{
		"clientConfig": ali.clientConfig,
	})
	if err != nil {
		return errors.Wrap(err, "创建配置客户端失败")
	}

	ali.client = client

	return nil
}

func (ali *aliACM) GetConfig(param vo.ConfigParam) (string, error) {
	// 获取配置
	return ali.client.GetConfig(param)
}

func (ali *aliACM) Listen(param vo.ConfigParam) error {
	return ali.client.ListenConfig(param)
}

func NewACM(opts acm.Options) (acm.ACM, error) {
	var ali aliACM

	// 设置默认值
	ali.validate = validator.New()
	ali.options = opts

	if err := ali.Validate(); err != nil {
		return nil, errors.Wrap(err, "参数验证失败")
	}

	if err := ali.configure(); err != nil {
		return nil, errors.Wrap(err, "ACM初始化失败")
	}

	return &ali, nil
}
