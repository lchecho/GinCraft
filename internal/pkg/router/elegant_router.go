package router

import "github.com/gin-gonic/gin"

// ElegantRouter 薄封装 gin.IRoutes，主要提供泛型 WrapRequestHandler 的路由入口。
type ElegantRouter struct {
	routes gin.IRoutes
}

// NewElegantRouter 从 *gin.Engine 构造
func NewElegantRouter(engine *gin.Engine) *ElegantRouter {
	return &ElegantRouter{routes: engine}
}

// NewElegantRouterGroup 从 *gin.RouterGroup 构造
func NewElegantRouterGroup(g *gin.RouterGroup) *ElegantRouter {
	return &ElegantRouter{routes: g}
}

// Group 创建子路由组（与 gin.RouterGroup.Group 语义一致）
func (r *ElegantRouter) Group(path string, mws ...gin.HandlerFunc) *ElegantRouter {
	type grouper interface {
		Group(string, ...gin.HandlerFunc) *gin.RouterGroup
	}
	g, ok := r.routes.(grouper)
	if !ok {
		return r
	}
	return NewElegantRouterGroup(g.Group(path, mws...))
}

// Use 添加中间件
func (r *ElegantRouter) Use(mws ...gin.HandlerFunc) gin.IRoutes {
	return r.routes.Use(mws...)
}

// combine 拷贝出新切片，避免共享 variadic 底层数组
func combine(mws []gin.HandlerFunc, h gin.HandlerFunc) []gin.HandlerFunc {
	out := make([]gin.HandlerFunc, 0, len(mws)+1)
	out = append(out, mws...)
	return append(out, h)
}

func (r *ElegantRouter) GET(p string, h gin.HandlerFunc, mws ...gin.HandlerFunc) gin.IRoutes {
	return r.routes.GET(p, combine(mws, h)...)
}

func (r *ElegantRouter) POST(p string, h gin.HandlerFunc, mws ...gin.HandlerFunc) gin.IRoutes {
	return r.routes.POST(p, combine(mws, h)...)
}

func (r *ElegantRouter) PUT(p string, h gin.HandlerFunc, mws ...gin.HandlerFunc) gin.IRoutes {
	return r.routes.PUT(p, combine(mws, h)...)
}

func (r *ElegantRouter) DELETE(p string, h gin.HandlerFunc, mws ...gin.HandlerFunc) gin.IRoutes {
	return r.routes.DELETE(p, combine(mws, h)...)
}

func (r *ElegantRouter) PATCH(p string, h gin.HandlerFunc, mws ...gin.HandlerFunc) gin.IRoutes {
	return r.routes.PATCH(p, combine(mws, h)...)
}
