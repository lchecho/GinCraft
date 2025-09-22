// @title GinCraft API
// @version 1.0
// @description GinCraft是一个基于Gin的Go Web框架，提供用户管理等功能
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 请输入Bearer token，例如：Bearer {token}

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/liuchen/gin-craft/internal/app"
	"github.com/liuchen/gin-craft/internal/pkg/config"
	"github.com/liuchen/gin-craft/internal/router"
	"github.com/liuchen/gin-craft/pkg/logger"
	"go.uber.org/zap"

	_ "github.com/liuchen/gin-craft/docs" // swagger文档
)

func main() {
	// 初始化应用
	if err := app.Init("config/config.yaml"); err != nil {
		fmt.Printf("Failed to initialize application: %v\n", err)
		os.Exit(1)
	}
	defer app.Close()

	// 设置Gin模式
	gin.SetMode(config.Config.App.Mode)

	// 创建路由
	r := router.InitRouter()

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.Config.App.Port),
		Handler:        r,
		ReadTimeout:    time.Duration(config.Config.App.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.Config.App.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// 启动HTTP服务器
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	logger.Info("Server started", zap.Int("port", config.Config.App.Port))

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// 设置关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exiting")
}
