package gin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

var pool = sync.Pool{New: func() interface{} {
	return new(Context)
}}

type Context struct {
	*gin.Context
}

// 扩展方法
func (c *Context) Success(data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func (c *Context) reset() {
	c.Context = nil
}

// 包装函数
type WrapHandlerFunc func(*Context)

// 包装
func Wrap(h WrapHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		context := pool.Get().(*Context) // 从缓冲池拿Context对象
		context.Context = c
		h(context)

		// 重置context
		context.reset()
		pool.Put(context) //  把Context对象放回缓冲池
	}
}
