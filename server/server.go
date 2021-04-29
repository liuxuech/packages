package server

import "github.com/gin-gonic/gin"

type Server struct {
	app *gin.Engine
}

func New() *Server {
	return &Server{app: gin.New()}
}

func (s *Server) Run(addr string) error {
	return s.app.Run(addr)
}

func (s *Server) Use(middleware ...HandlerFunc) gin.IRoutes {
	s.Use(middleware...)
	return s.app
}

func (s *Server) GET(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	s.GET(relativePath, handlers...)
	return s.app
}

func (s *Server) POST(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	s.GET(relativePath, handlers...)
	return s.app
}

func (s *Server) PUT(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	s.GET(relativePath, handlers...)
	return s.app
}

func (s *Server) DELETE(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	s.GET(relativePath, handlers...)
	return s.app
}
