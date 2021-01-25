package acm

import (
	"github.com/go-playground/validator/v10"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/pkg/errors"
)

type nacosACM struct {
	options  Options
	validate *validator.Validate

	client       config_client.IConfigClient
	clientConfig constant.ClientConfig
	serverConfig constant.ServerConfig
}

func (nacos *nacosACM) Listen(param vo.ConfigParam) error {
	return nacos.client.ListenConfig(param)
}

func (nacos *nacosACM) GetConfig(param vo.ConfigParam) (string, error) {
	return nacos.client.GetConfig(param)
}

func (nacos *nacosACM) Validate() error {
	return nacos.validate.Struct(nacos.options)
}

func (nacos *nacosACM) configure() error {
	nacos.clientConfig = constant.ClientConfig{
		TimeoutMs:   5 * 1000,
		NamespaceId: nacos.options.Namespace,
	}

	nacos.serverConfig = constant.ServerConfig{
		IpAddr: nacos.options.Host,
		Port:   nacos.options.Port,
	}

	c, err := clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  &nacos.clientConfig,
		ServerConfigs: []constant.ServerConfig{nacos.serverConfig},
	})
	if err != nil {
		return errors.Wrap(err, "新建配置客户端错误")
	}
	nacos.client = c

	return nil
}

func newACM(opts Options) (ACM, error) {
	var nacos nacosACM

	// 设置默认值
	nacos.validate = validator.New()
	nacos.options = opts

	// 参数验证
	if err := nacos.Validate(); err != nil {
		return nil, errors.Wrap(err, "参数验证失败")
	}

	if err := nacos.configure(); err != nil {
		return nil, errors.Wrap(err, "ACM初始化失败")
	}

	return &nacos, nil
}
