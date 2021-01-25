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

type aliACM struct {
	options acm.Options

	client config_client.IConfigClient

	clientConfig constant.ClientConfig

	validate *validator.Validate
}

func (ali *aliACM) Init() error {
	access, ok := ali.options.Others[accessKey]
	if !ok || access == "" {
		return errors.New("缺少必要参数: accessKey")
	}
	secret, ok := ali.options.Others[secretKey]
	if !ok || secret == "" {
		return errors.New("缺少必要参数: secretKey")
	}
	if ali.options.NamespaceId == "" {
		return errors.New("缺少必要参数: NamespaceId")
	}

	// access和secret 在设置的时候就指定了类型为string
	ali.clientConfig = constant.ClientConfig{
		AccessKey:   access.(string),
		SecretKey:   secret.(string),
		TimeoutMs:   5 * 1000,
		NamespaceId: ali.options.NamespaceId,
		Endpoint:    fmt.Sprintf("%s:%d", ali.options.Host, ali.options.Port),
	}
	if ali.options.CacheDir != "" {
		ali.clientConfig.CacheDir = ali.options.CacheDir
	}
	if ali.options.LogDir != "" {
		ali.clientConfig.LogDir = ali.options.LogDir
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

func (ali *aliACM) GetConfig() (string, error) {
	// 获取配置
	content, err := ali.client.GetConfig(vo.ConfigParam{
		DataId: "beaver.yaml",
		Group:  "DEFAULT_GROUP"})
	if err != nil {
		return "", err
	}

	return content, nil
}

func NewACM(opts ...acm.Option) acm.ACM {
	var ali aliACM

	// 设置默认值
	ali.validate = validator.New()
	ali.options.Host = "127.0.0.1"
	ali.options.Port = 8848
	ali.options.Others = make(map[string]interface{})

	for _, o := range opts {
		o(&ali.options)
	}

	return &ali
}
