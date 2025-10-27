package router

import (
	"github.com/gin-gonic/gin"
	"github.com/liuchen/gin-craft/internal/controller"
	"github.com/liuchen/gin-craft/internal/middleware"
	elegantRouter "github.com/liuchen/gin-craft/internal/pkg/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRouter 初始化路由（使用struct controller）
func InitRouter() *gin.Engine {
	r := gin.New()

	// 使用全局中间件
	r.Use(middleware.ContextMiddleware()) // 自定义Context中间件（必须在最前面）
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery()) // 使用恢复中间件
	r.Use(middleware.Cors())

	// 创建优雅路由器
	elegantR := elegantRouter.NewElegantRouter(r)

	// 健康检查（无中间件）
	elegantR.GET("/health", elegantRouter.WrapHandler(func(c *gin.Context) (interface{}, error) {
		return gin.H{"status": "ok"}, nil
	}))

	// Swagger文档路由
	elegantR.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API 路由组
	api := elegantR.Group("/api")
	{
		// v1 版本
		v1 := api.Group("/v1")
		{
			// 公开的用户路由（无中间件）
			// 使用统一的参数解析处理器
			// 用户路由
			userController := controller.NewUserController()
			v1.POST("/user/register", elegantRouter.WrapRequestHandler(userController.Register))
			v1.POST("/user/login", elegantRouter.WrapRequestHandler(userController.Login))
			v1.POST("/user/list", elegantRouter.WrapRequestHandler(userController.List))
			v1.GET("/user/info", elegantRouter.WrapRequestHandler(userController.Info))
			v1.POST("/user/edit", elegantRouter.WrapRequestHandler(userController.Update))
			v1.POST("/user/delete", elegantRouter.WrapRequestHandler(userController.Delete))

			// 需要认证的用户路由（使用中间件）
			authUser := v1.Group("/user", middleware.AuthMiddleware())
			{
				authUser.GET("/info", elegantRouter.WrapRequestHandler(userController.Info))
				// 可以为单个路由添加额外中间件
				authUser.GET("/profile", elegantRouter.WrapRequestHandler(userController.Info), middleware.RateLimitMiddleware())
			}

			// 管理员路由（多个中间件）
			admin := v1.Group("/admin", middleware.AuthMiddleware(), middleware.AdminAuthMiddleware())
			{
				admin.GET("/users", elegantRouter.WrapHandler(func(c *gin.Context) (interface{}, error) {
					return gin.H{"message": "管理员用户列表"}, nil
				}))
			}

			// API路由（使用API密钥认证）
			apiRoutes := v1.Group("/api", middleware.ValidateAPIKeyMiddleware())
			{
				apiRoutes.GET("/data", elegantRouter.WrapHandler(func(c *gin.Context) (interface{}, error) {
					return gin.H{"data": "API数据"}, nil
				}))
			}
		}
	}

	return r
}
