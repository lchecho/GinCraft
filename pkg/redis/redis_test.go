package redis

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestRedisClient 测试Redis客户端连接
func TestRedisClient(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
		PoolSize: 10,
	}

	client := NewClient(config)
	assert.NotNil(t, client)

	// 注意：这个测试需要本地运行Redis服务
	// 如果没有Redis服务，测试将失败
	err := client.Connect()
	if err != nil {
		t.Skipf("Redis not available: %v", err)
		return
	}
	defer client.Close()

	// 测试Ping
	err = client.Ping()
	assert.NoError(t, err)
}

// TestRedisOperations 测试Redis基本操作
func TestRedisOperations(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
		PoolSize: 10,
	}

	client := NewClient(config)
	err := client.Connect()
	if err != nil {
		t.Skipf("Redis not available: %v", err)
		return
	}
	defer client.Close()

	ctx := context.Background()
	testKey := "test:key"
	testValue := "test_value"

	// 清理测试数据
	defer client.Del(ctx, testKey)

	// 测试Set操作
	err = client.Set(ctx, testKey, testValue, 10*time.Second)
	assert.NoError(t, err)

	// 测试Get操作
	value, err := client.Get(ctx, testKey)
	assert.NoError(t, err)
	assert.Equal(t, testValue, value)

	// 测试Exists操作
	count, err := client.Exists(ctx, testKey)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	// 测试TTL操作
	ttl, err := client.TTL(ctx, testKey)
	assert.NoError(t, err)
	assert.Greater(t, ttl.Seconds(), float64(0))

	// 测试Expire操作
	err = client.Expire(ctx, testKey, 20*time.Second)
	assert.NoError(t, err)

	// 测试Del操作
	err = client.Del(ctx, testKey)
	assert.NoError(t, err)

	// 验证删除后不存在
	count, err = client.Exists(ctx, testKey)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

// TestRedisExample 测试Redis示例
func TestRedisExample(t *testing.T) {
	config := &Config{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
		PoolSize: 10,
	}

	client := NewClient(config)
	err := client.Connect()
	if err != nil {
		t.Skipf("Redis not available: %v", err)
		return
	}
	defer client.Close()

	example := NewExample(client)
	ctx := context.Background()
	testKey := "example:key"

	// 清理测试数据
	defer example.DeleteCache(ctx, testKey)

	// 测试设置缓存
	err = example.SetCache(ctx, testKey, "example_value", 10*time.Second)
	assert.NoError(t, err)

	// 测试获取缓存
	value, err := example.GetCache(ctx, testKey)
	assert.NoError(t, err)
	assert.Equal(t, "example_value", value)

	// 测试检查缓存存在
	exists, err := example.CheckCacheExists(ctx, testKey)
	assert.NoError(t, err)
	assert.True(t, exists)

	// 测试获取TTL
	ttl, err := example.GetCacheTTL(ctx, testKey)
	assert.NoError(t, err)
	assert.Greater(t, ttl.Seconds(), float64(0))

	// 测试删除缓存
	err = example.DeleteCache(ctx, testKey)
	assert.NoError(t, err)

	// 验证删除后不存在
	exists, err = example.CheckCacheExists(ctx, testKey)
	assert.NoError(t, err)
	assert.False(t, exists)
}
