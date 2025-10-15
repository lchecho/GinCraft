package redis

import (
	"context"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config Redis配置
type Config struct {
	Host         string
	Port         int
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	MaxRetries   int
	DialTimeout  int // 连接超时时间(秒)
	ReadTimeout  int // 读取超时时间(秒)
	WriteTimeout int // 写入超时时间(秒)
}

// Client Redis数据库实现
type Client struct {
	client *redis.Client
	config *Config
	mu     sync.RWMutex
}

// NewClient 创建Redis客户端实例
func NewClient(config *Config) *Client {
	// 设置默认值
	if config.PoolSize == 0 {
		config.PoolSize = 10
	}
	if config.MinIdleConns == 0 {
		config.MinIdleConns = 5
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.DialTimeout == 0 {
		config.DialTimeout = 5
	}
	if config.ReadTimeout == 0 {
		config.ReadTimeout = 3
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = 3
	}

	return &Client{
		config: config,
	}
}

// Connect 连接Redis
func (r *Client) Connect() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.client != nil {
		return nil // 已经连接
	}

	addr := fmt.Sprintf("%s:%d", r.config.Host, r.config.Port)

	r.client = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     r.config.Password,
		DB:           r.config.DB,
		PoolSize:     r.config.PoolSize,
		MinIdleConns: r.config.MinIdleConns,
		MaxRetries:   r.config.MaxRetries,
		DialTimeout:  time.Duration(r.config.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(r.config.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(r.config.WriteTimeout) * time.Second,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return nil
}

// GetClient 获取Redis客户端
func (r *Client) GetClient() *redis.Client {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.client
}

// Close 关闭Redis连接
func (r *Client) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.client != nil {
		if err := r.client.Close(); err != nil {
			return fmt.Errorf("failed to close Redis connection: %w", err)
		}
		r.client = nil
	}
	return nil
}

// Ping 测试Redis连接
func (r *Client) Ping() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.client == nil {
		return fmt.Errorf("redis not connected")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return r.client.Ping(ctx).Err()
}

// Set 设置键值
func (r *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.client == nil {
		return fmt.Errorf("redis not connected")
	}

	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get 获取键值
func (r *Client) Get(ctx context.Context, key string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.client == nil {
		return "", fmt.Errorf("redis not connected")
	}

	return r.client.Get(ctx, key).Result()
}

// Del 删除键
func (r *Client) Del(ctx context.Context, keys ...string) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.client == nil {
		return fmt.Errorf("redis not connected")
	}

	return r.client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func (r *Client) Exists(ctx context.Context, keys ...string) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.client == nil {
		return 0, fmt.Errorf("redis not connected")
	}

	return r.client.Exists(ctx, keys...).Result()
}

// Expire 设置过期时间
func (r *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.client == nil {
		return fmt.Errorf("redis not connected")
	}

	return r.client.Expire(ctx, key, expiration).Err()
}

// TTL 获取过期时间
func (r *Client) TTL(ctx context.Context, key string) (time.Duration, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.client == nil {
		return 0, fmt.Errorf("redis not connected")
	}

	return r.client.TTL(ctx, key).Result()
}

// SetJSON 设置JSON数据到Redis
func (r *Client) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := jsoniter.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return r.client.Set(ctx, key, data, expiration).Err()
}

// GetJSON 从Redis获取JSON数据
func (r *Client) GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := r.client.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return fmt.Errorf("failed to get from Redis: %w", err)
	}

	if data == "" {
		return nil
	}

	if err = jsoniter.Unmarshal([]byte(data), dest); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

// SetWithNX 仅当键不存在时设置（分布式锁）
func (r *Client) SetWithNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, expiration).Result()
}

// Increment 递增计数器
func (r *Client) Increment(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

// IncrementBy 递增指定值
func (r *Client) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
	return r.client.IncrBy(ctx, key, value).Result()
}

// Decrement 递减计数器
func (r *Client) Decrement(ctx context.Context, key string) (int64, error) {
	return r.client.Decr(ctx, key).Result()
}

// SetExpire 设置过期时间
func (r *Client) SetExpire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}

// GetTTL 获取过期时间
func (r *Client) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

// HSet 设置哈希表字段
func (r *Client) HSet(ctx context.Context, key string, field string, value interface{}) error {
	return r.client.HSet(ctx, key, field, value).Err()
}

// HGet 获取哈希表字段
func (r *Client) HGet(ctx context.Context, key string, field string) (string, error) {
	return r.client.HGet(ctx, key, field).Result()
}

// HGetAll 获取哈希表所有字段
func (r *Client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.client.HGetAll(ctx, key).Result()
}

// HDel 删除哈希表字段
func (r *Client) HDel(ctx context.Context, key string, fields ...string) error {
	return r.client.HDel(ctx, key, fields...).Err()
}

// LPush 从左侧推入列表
func (r *Client) LPush(ctx context.Context, key string, values ...interface{}) error {
	return r.client.LPush(ctx, key, values...).Err()
}

// RPush 从右侧推入列表
func (r *Client) RPush(ctx context.Context, key string, values ...interface{}) error {
	return r.client.RPush(ctx, key, values...).Err()
}

// LPop 从左侧弹出列表
func (r *Client) LPop(ctx context.Context, key string) (string, error) {
	return r.client.LPop(ctx, key).Result()
}

// RPop 从右侧弹出列表
func (r *Client) RPop(ctx context.Context, key string) (string, error) {
	return r.client.RPop(ctx, key).Result()
}

// LLen 获取列表长度
func (r *Client) LLen(ctx context.Context, key string) (int64, error) {
	return r.client.LLen(ctx, key).Result()
}

// SAdd 添加集合成员
func (r *Client) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return r.client.SAdd(ctx, key, members...).Err()
}

// SRem 删除集合成员
func (r *Client) SRem(ctx context.Context, key string, members ...interface{}) error {
	return r.client.SRem(ctx, key, members...).Err()
}

// SMembers 获取集合所有成员
func (r *Client) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.client.SMembers(ctx, key).Result()
}

// SIsMember 检查是否是集合成员
func (r *Client) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return r.client.SIsMember(ctx, key, member).Result()
}

// ZAdd 添加有序集合成员
func (r *Client) ZAdd(ctx context.Context, key string, score float64, member interface{}) error {
	return r.client.ZAdd(ctx, key, redis.Z{Score: score, Member: member}).Err()
}

// ZRem 删除有序集合成员
func (r *Client) ZRem(ctx context.Context, key string, members ...interface{}) error {
	return r.client.ZRem(ctx, key, members...).Err()
}

// ZRange 获取有序集合范围成员
func (r *Client) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.client.ZRange(ctx, key, start, stop).Result()
}

// ZRangeWithScores 获取有序集合范围成员（包含分数）
func (r *Client) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	return r.client.ZRangeWithScores(ctx, key, start, stop).Result()
}

// ZCard 获取有序集合成员数量
func (r *Client) ZCard(ctx context.Context, key string) (int64, error) {
	return r.client.ZCard(ctx, key).Result()
}

// ZRevRange 获取有序集合范围成员（倒序）
func (r *Client) ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.client.ZRevRange(ctx, key, start, stop).Result()
}

// ZRevRangeWithScores 获取有序集合范围成员（倒序，包含分数）
func (r *Client) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	return r.client.ZRevRangeWithScores(ctx, key, start, stop).Result()
}

// ZScore 获取有序集合成员分数
func (r *Client) ZScore(ctx context.Context, key string, member string) (float64, error) {
	return r.client.ZScore(ctx, key, member).Result()
}
