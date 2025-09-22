package router

import (
	"github.com/gin-gonic/gin"
	"github.com/liuchen/gin-craft/pkg/response"
)

// HandlerFunc 统一的处理器函数类型
type HandlerFunc func(ctx *gin.Context) (interface{}, error)

// ElegantRouter 优雅的路由器，封装 gin.IRoutes
type ElegantRouter struct {
	routes gin.IRoutes
	engine *gin.Engine
	group  *gin.RouterGroup
}

// NewElegantRouter 创建优雅路由器（从 Engine）
func NewElegantRouter(engine *gin.Engine) *ElegantRouter {
	return &ElegantRouter{routes: engine, engine: engine}
}

// NewElegantRouterGroup 创建优雅路由器（从 RouterGroup）
func NewElegantRouterGroup(group *gin.RouterGroup) *ElegantRouter {
	return &ElegantRouter{routes: group, group: group}
}

// Group 创建路由组，支持可选中间件
func (r *ElegantRouter) Group(relativePath string, middleware ...gin.HandlerFunc) *ElegantRouter {
	var group *ElegantRouter

	if r.engine != nil {
		group = NewElegantRouterGroup(r.engine.Group(relativePath))
	} else if r.group != nil {
		group = NewElegantRouterGroup(r.group.Group(relativePath))
	} else {
		return r
	}

	// 如果提供了中间件，则添加到路由组
	if len(middleware) > 0 {
		group.Use(middleware...)
	}

	return group
}

// Use 添加中间件
func (r *ElegantRouter) Use(middleware ...gin.HandlerFunc) gin.IRoutes {
	return r.routes.Use(middleware...)
}

// WithMiddleware 为当前路由器添加中间件并返回新的路由器实例
func (r *ElegantRouter) WithMiddleware(middleware ...gin.HandlerFunc) *ElegantRouter {
	newRouter := &ElegantRouter{
		routes: r.routes,
		engine: r.engine,
		group:  r.group,
	}
	newRouter.Use(middleware...)
	return newRouter
}

// GET 处理 GET 请求
func (r *ElegantRouter) GET(relativePath string, handlerFunc HandlerFunc, middleware ...gin.HandlerFunc) gin.IRoutes {
	handlers := append(middleware, r.wrapHandler(handlerFunc))
	return r.routes.GET(relativePath, handlers...)
}

// POST 处理 POST 请求
func (r *ElegantRouter) POST(relativePath string, handlerFunc HandlerFunc, middleware ...gin.HandlerFunc) gin.IRoutes {
	handlers := append(middleware, r.wrapHandler(handlerFunc))
	return r.routes.POST(relativePath, handlers...)
}

// PUT 处理 PUT 请求
func (r *ElegantRouter) PUT(relativePath string, handlerFunc HandlerFunc, middleware ...gin.HandlerFunc) gin.IRoutes {
	handlers := append(middleware, r.wrapHandler(handlerFunc))
	return r.routes.PUT(relativePath, handlers...)
}

// DELETE 处理 DELETE 请求
func (r *ElegantRouter) DELETE(relativePath string, handlerFunc HandlerFunc, middleware ...gin.HandlerFunc) gin.IRoutes {
	handlers := append(middleware, r.wrapHandler(handlerFunc))
	return r.routes.DELETE(relativePath, handlers...)
}

// PATCH 处理 PATCH 请求
func (r *ElegantRouter) PATCH(relativePath string, handlerFunc HandlerFunc, middleware ...gin.HandlerFunc) gin.IRoutes {
	handlers := append(middleware, r.wrapHandler(handlerFunc))
	return r.routes.PATCH(relativePath, handlers...)
}

// wrapHandler 通用包装处理函数
func (r *ElegantRouter) wrapHandler(handlerFunc HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := handlerFunc(c)
		r.handleResponse(c, data, err)
	}
}

// handleResponse 统一处理响应
func (r *ElegantRouter) handleResponse(c *gin.Context, data interface{}, err error) {
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, data)
}
