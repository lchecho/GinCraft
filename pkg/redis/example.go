package redis

import (
	"context"
	"fmt"
	"time"
)

// Example Redis使用示例
type Example struct {
	client *Client
}

// NewExample 创建Redis示例实例
func NewExample(client *Client) *Example {
	return &Example{
		client: client,
	}
}

// SetCache 设置缓存
func (r *Example) SetCache(ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration)
}

// GetCache 获取缓存
func (r *Example) GetCache(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key)
}

// DeleteCache 删除缓存
func (r *Example) DeleteCache(ctx context.Context, key string) error {
	return r.client.Del(ctx, key)
}

// CheckCacheExists 检查缓存是否存在
func (r *Example) CheckCacheExists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetCacheTTL 获取缓存过期时间
func (r *Example) GetCacheTTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key)
}

// SetCacheWithRefresh 设置缓存并支持刷新过期时间
func (r *Example) SetCacheWithRefresh(ctx context.Context, key string, value string, expiration time.Duration) error {
	// 先尝试获取现有值
	exists, err := r.CheckCacheExists(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to check cache exists: %w", err)
	}

	if exists {
		// 如果已存在，更新过期时间
		if err := r.client.Expire(ctx, key, expiration); err != nil {
			return fmt.Errorf("failed to refresh cache expiration: %w", err)
		}
	} else {
		// 如果不存在，设置新值
		if err := r.client.Set(ctx, key, value, expiration); err != nil {
			return fmt.Errorf("failed to set cache: %w", err)
		}
	}

	return nil
}
