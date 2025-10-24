package middleware

import (
	"github.com/liuchen/gin-craft/internal/pkg/context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/liuchen/gin-craft/internal/constant"
	"github.com/liuchen/gin-craft/pkg/errors"
	"github.com/liuchen/gin-craft/pkg/response"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, errors.New(constant.Unauthorized, "缺少认证令牌"))
			c.Abort()
			return
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(c, errors.New(constant.Unauthorized, "认证令牌格式错误"))
			c.Abort()
			return
		}

		// 提取token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			response.Error(c, errors.New(constant.Unauthorized, "认证令牌为空"))
			c.Abort()
			return
		}

		// 使用应用 Context 记录与传播认证信息
		appCtx := context.MustGetContext(c)
		appCtx.SetCustomField("token_length", len(token))
		appCtx.SetCustomField("token", token)
		// TODO: 在接入真实 JWT 后，从 claims 中解析用户信息
		appCtx.SetUser("123", "user_123", "user")

		c.Next()
	}
}

// AdminAuthMiddleware 管理员认证中间件
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先检查是否已经通过基础认证
		var token string
		if v, ok := context.MustGetContext(c).GetCustomField("token"); ok {
			token, _ = v.(string)
		}
		if token == "" {
			response.Error(c, errors.New(constant.Unauthorized, "需要先进行认证"))
			c.Abort()
			return
		}

		// 这里可以添加管理员权限检查逻辑
		// 简化处理，假设token中包含admin标识
		if !strings.Contains(token, "admin") {
			response.Error(c, errors.New(constant.Forbidden, "需要管理员权限"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这里可以添加限流逻辑
		// 简化处理，直接通过
		c.Next()
	}
}

// ValidateAPIKeyMiddleware API密钥验证中间件
func ValidateAPIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			response.Error(c, errors.New(constant.Unauthorized, "缺少API密钥"))
			c.Abort()
			return
		}

		// 这里可以添加API密钥验证逻辑
		// 简化处理，检查是否为预设值
		if apiKey != "your-api-key" {
			response.Error(c, errors.New(constant.Unauthorized, "无效的API密钥"))
			c.Abort()
			return
		}

		c.Next()
	}
}
