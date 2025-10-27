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

// InitRedis 初始化Redis连接
func InitRedis() error {
	var err error
	once.Do(func() {
		cfg := config.Config.Redis

		// 创建Redis客户端实例
		redisConfig := &pkgredis.Config{
			Host:         cfg.Host,
			Port:         cfg.Port,
			Password:     cfg.Password,
			DB:           cfg.DB,
			PoolSize:     cfg.PoolSize,
			MinIdleConns: 5,
			MaxRetries:   3,
			DialTimeout:  5,
			ReadTimeout:  3,
			WriteTimeout: 3,
		}

		client = pkgredis.NewClient(redisConfig)
		err = client.Connect()
		if err != nil {
			return
		}

		// 测试连接
		err = client.Ping()
		if err != nil {
			return
		}
	})
	return err
}

// GetClient 获取Redis客户端
func GetClient() *redis.Client {
	if client == nil {
		return nil
	}
	return client.GetClient()
}

// GetRedisClient 获取封装的Client接口
func GetRedisClient() *pkgredis.Client {
	return client
}

// Close 关闭Redis连接
func Close() {
	if client != nil {
		_ = client.Close()
	}
}
