package redis

import (
	"sync"

	"github.com/liuchen/gin-craft/internal/pkg/config"
	pkgredis "github.com/liuchen/gin-craft/pkg/redis"
	"github.com/redis/go-redis/v9"
)

var (
	client *pkgredis.Client
	once   sync.Once
)

// InitRedis 初始化 Redis 连接
func InitRedis() error {
	var err error
	once.Do(func() {
		cfg := config.Config.Redis
		redisConfig := &pkgredis.Config{
			Host:         cfg.Host,
			Port:         cfg.Port,
			Password:     cfg.Password,
			DB:           cfg.DB,
			PoolSize:     cfg.PoolSize,
			MinIdleConns: cfg.MinIdleConns,
			MaxRetries:   cfg.MaxRetries,
			DialTimeout:  cfg.DialTimeout,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		}

		client = pkgredis.NewClient(redisConfig)
		if err = client.Connect(); err != nil {
			return
		}
		err = client.Ping()
	})
	return err
}

// GetClient 获取底层 *redis.Client
func GetClient() *redis.Client {
	if client == nil {
		return nil
	}
	return client.GetClient()
}

// GetRedisClient 获取封装后的 Client
func GetRedisClient() *pkgredis.Client {
	return client
}

// Close 关闭 Redis 连接
func Close() {
	if client != nil {
		_ = client.Close()
	}
}
