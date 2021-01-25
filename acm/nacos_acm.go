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
	options Options

	client config_client.IConfigClient

	clientCfg constant.ClientConfig
	serverCfg constant.ServerConfig

	validate *validator.Validate
}

func (nacos *nacosACM) GetConfig() (string, error) {
	content, err := nacos.client.GetConfig(vo.ConfigParam{
		DataId: "clear.yaml",
		Group:  "DEFAULT_GROUP",
	})
	return content, err
}

func (nacos *nacosACM) Init() error {
	// namespaceId 没有默认值，但是又是必须的值，所以需要验证
	if err := nacos.validate.Var(nacos.options.NamespaceId, "required"); err != nil {
		return errors.Wrap(err, "缺少参数")
	}

	nacos.clientCfg = constant.ClientConfig{
		TimeoutMs:   5 * 1000,
		NamespaceId: nacos.options.NamespaceId,
	}

	nacos.serverCfg = constant.ServerConfig{
		IpAddr: nacos.options.IpAddr,
		Port:   nacos.options.Port,
	}

	c, err := clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  &nacos.clientCfg,
		ServerConfigs: []constant.ServerConfig{nacos.serverCfg},
	})
	if err != nil {
		return errors.Wrap(err, "新建配置客户端错误")
	}

	nacos.client = c

	return nil
}

func newACM(opts ...Option) ACM {
	var nacos nacosACM

	// 设置默认值
	nacos.validate = validator.New()
	nacos.options.TimeoutMs = 5 * 1000
	nacos.options.IpAddr = "127.0.0.1"
	nacos.options.Port = 8848

	for _, o := range opts {
		o(&nacos.options)
	}

	return &nacos
}
