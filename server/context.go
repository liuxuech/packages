package server

import (
	"github.com/gin-gonic/gin"
)

type HandlerFunc func(*Context)

type Context struct {
	*gin.Context
}

func (ctx *Context) Success(code int, data interface{}) {
	ctx.JSON(code, data)
}
