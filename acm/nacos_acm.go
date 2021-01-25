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

	clientConfig constant.ClientConfig
	serverConfig constant.ServerConfig

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
	if nacos.options.NamespaceId == "" {
		return errors.New("缺少必要参数: NamespaceId")
	}

	nacos.clientConfig = constant.ClientConfig{
		TimeoutMs:   5 * 1000,
		NamespaceId: nacos.options.NamespaceId,
	}

	if nacos.options.CacheDir != "" {
		nacos.clientConfig.CacheDir = nacos.options.CacheDir
	}
	if nacos.options.LogDir != "" {
		nacos.clientConfig.LogDir = nacos.options.LogDir
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

func newACM(opts ...Option) ACM {
	var nacos nacosACM

	// 设置默认值
	nacos.validate = validator.New()
	nacos.options.Host = "127.0.0.1"
	nacos.options.Port = 8848
	nacos.options.Others = make(map[string]interface{})

	for _, o := range opts {
		o(&nacos.options)
	}

	return &nacos
}
