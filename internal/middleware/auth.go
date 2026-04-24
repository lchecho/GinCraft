package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/liuchen/gin-craft/internal/constant"
	appctx "github.com/liuchen/gin-craft/internal/pkg/context"
	"github.com/liuchen/gin-craft/internal/pkg/errors"
	"github.com/liuchen/gin-craft/internal/pkg/response"
	"github.com/liuchen/gin-craft/internal/pkg/config"
)

const (
	// ctxHasTokenKey 仅标记请求携带了 token，不存明文 token 本身
	ctxHasTokenKey = "has_token"
)

// AuthMiddleware Bearer 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			unauthorized(c, "认证令牌无效")
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			unauthorized(c, "认证令牌为空")
			return
		}

		// TODO: 接入真实 JWT 后，从 claims 中解析用户信息
		userID, username, role := "123", "user_123", "user"

		appCtx := appctx.MustGetContext(c)
		appCtx.SetUser(userID, username, role)
		appCtx.SetCustomField(ctxHasTokenKey, true) // 仅标记，不写 token 原文
		c.Next()
	}
}

// AdminAuthMiddleware 管理员认证中间件（必须在 AuthMiddleware 之后）
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appCtx := appctx.MustGetContext(c)
		if appCtx.GetUserRole() != "admin" {
			response.Error(c, errors.New(constant.Forbidden, "需要管理员权限"))
			c.Abort()
			return
		}
		c.Next()
	}
}

// RateLimitMiddleware 限流中间件占位实现
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 接入真实的基于 token/ip 的限流
		c.Next()
	}
}

// ValidateAPIKeyMiddleware API 密钥校验，密钥从 config.app.api_key 读取
func ValidateAPIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		expected := config.Config.App.APIKey
		apiKey := c.GetHeader("X-API-Key")
		if expected == "" || apiKey == "" || apiKey != expected {
			response.Error(c, errors.New(constant.Unauthorized, "无效的API密钥"))
			c.Abort()
			return
		}
		c.Next()
	}
}

func unauthorized(c *gin.Context, msg string) {
	response.Error(c, errors.New(constant.Unauthorized, msg))
	c.Abort()
}
