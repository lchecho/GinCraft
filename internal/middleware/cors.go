package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/liuchen/gin-craft/internal/pkg/config"
)

// Cors 跨域中间件（由 config.cors 配置驱动）
func Cors() gin.HandlerFunc {
	cfg := config.Config.CORS

	methods := strings.Join(defaultIfEmpty(cfg.AllowMethods,
		[]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}), ", ")
	headers := strings.Join(defaultIfEmpty(cfg.AllowHeaders,
		[]string{"Content-Type", "Authorization", "X-Requested-With", "X-API-Key"}), ", ")

	allowed := make(map[string]struct{}, len(cfg.AllowOrigins))
	for _, o := range cfg.AllowOrigins {
		allowed[o] = struct{}{}
	}
	_, allowAny := allowed["*"]
	allowAll := len(cfg.AllowOrigins) == 0 || allowAny

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		switch {
		case allowAll && !cfg.AllowCredentials:
			// 不带 credentials 的全开模式
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		case origin != "" && (allowAll || hasKey(allowed, origin)):
			// 指定来源，回显 Origin（这是 credentials=true 的唯一合法形式）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Vary", "Origin")
		}

		if cfg.AllowCredentials {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		c.Writer.Header().Set("Access-Control-Allow-Methods", methods)
		c.Writer.Header().Set("Access-Control-Allow-Headers", headers)
		if cfg.MaxAge > 0 {
			c.Writer.Header().Set("Access-Control-Max-Age", strconv.Itoa(cfg.MaxAge))
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func defaultIfEmpty(v, d []string) []string {
	if len(v) == 0 {
		return d
	}
	return v
}

func hasKey(set map[string]struct{}, k string) bool {
	_, ok := set[k]
	return ok
}
