package web

import "github.com/gin-gonic/gin"

type Web struct {
	*gin.Engine
}

func (w *Web) Use(middleware ...HandlerFunc) {
	w.Use(middleware...)
}

func (w *Web) GET(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	w.GET(relativePath, handlers...)
	return w.Engine
}

func (w *Web) POST(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	w.GET(relativePath, handlers...)
	return w.Engine
}

func (w *Web) PUT(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	w.GET(relativePath, handlers...)
	return w.Engine
}

func (w *Web) DELETE(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	w.GET(relativePath, handlers...)
	return w.Engine
}
