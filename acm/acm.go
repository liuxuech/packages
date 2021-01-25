package acm

import (
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type ACM interface {
	GetConfig(vo.ConfigParam) (string, error)
	Listen(vo.ConfigParam) error
}

type ListenFunc func(namespace, group, dataId, data string)

// 公共的配置参数
type Options struct {
	Host      string `validate:"required"`
	Port      uint64 `validate:"required"`
	Namespace string `validate:"required"` // 命名空间，ep：dev、test、prod

	// 其他实现所需的配置可以存放在这里
	Others map[string]interface{}
}

func NewACM(opts Options) (ACM, error) {
	return newACM(opts)
}
