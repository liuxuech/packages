package web

import "github.com/gin-gonic/gin"

type HandlerFunc func(*Context)

type Context struct {
	*gin.Context
}
