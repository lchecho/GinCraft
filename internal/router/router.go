package router

import (
	"github.com/gin-gonic/gin"
	"github.com/liuchen/gin-craft/internal/controller"
	"github.com/liuchen/gin-craft/internal/middleware"
	er "github.com/liuchen/gin-craft/internal/pkg/router"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(
		middleware.ContextMiddleware(),
		middleware.Logger(),
		middleware.Recovery(),
		middleware.Cors(),
	)

	elegantR := er.NewElegantRouter(r)

	elegantR.GET("/health", er.WrapHandler(func(c *gin.Context) (interface{}, error) {
		return gin.H{"status": "ok"}, nil
	}))
	elegantR.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userCtrl := controller.NewUserController()

	api := elegantR.Group("/api")
	v1 := api.Group("/v1")

	// 公开路由（无认证）
	publicUser := v1.Group("/user")
	{
		publicUser.POST("/register", er.WrapRequestHandler(userCtrl.Register))
		publicUser.POST("/login", er.WrapRequestHandler(userCtrl.Login))
	}

	// 认证路由
	authUser := v1.Group("/user", middleware.AuthMiddleware())
	{
		authUser.GET("/info", er.WrapRequestHandler(userCtrl.Info))
		authUser.POST("/list", er.WrapRequestHandler(userCtrl.List))
		authUser.POST("/edit", er.WrapRequestHandler(userCtrl.Update))
		authUser.POST("/delete", er.WrapRequestHandler(userCtrl.Delete))
		authUser.GET("/profile", er.WrapRequestHandler(userCtrl.Info), middleware.RateLimitMiddleware())
	}

	admin := v1.Group("/admin", middleware.AuthMiddleware(), middleware.AdminAuthMiddleware())
	{
		admin.GET("/users", er.WrapHandler(func(c *gin.Context) (interface{}, error) {
			return gin.H{"message": "管理员用户列表"}, nil
		}))
	}

	apiRoutes := v1.Group("/api", middleware.ValidateAPIKeyMiddleware())
	{
		apiRoutes.GET("/data", er.WrapHandler(func(c *gin.Context) (interface{}, error) {
			return gin.H{"data": "API数据"}, nil
		}))
	}

	return r
}
