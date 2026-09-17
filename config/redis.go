package config

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

type mockEntry struct {
	val    string
	expiry time.Time
}

var (
	mockStore = struct {
		sync.RWMutex
		m map[string]mockEntry
	}{m: make(map[string]mockEntry)}
)

func InitRedis() {
	// 初始化 redis 连接
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: 20,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		panic("Redis连接失败: " + err.Error())
	}
}

// RedisSet 在有真实 Redis 时调用 Redis，否则使用内存替代（方便测试）
func RedisSet(ctx context.Context, key string, val string, expiration time.Duration) error {
	if RedisClient != nil {
		return RedisClient.Set(ctx, key, val, expiration).Err()
	}
	mockStore.Lock()
	defer mockStore.Unlock()
	e := mockEntry{val: val}
	if expiration > 0 {
		e.expiry = time.Now().Add(expiration)
	}
	mockStore.m[key] = e
	return nil
}

// RedisGetDel 在有真实 Redis 时调用 Redis 的 GetDel，否则使用内存替代
func RedisGetDel(ctx context.Context, key string) (string, error) {
	if RedisClient != nil {
		return RedisClient.GetDel(ctx, key).Result()
	}
	mockStore.Lock()
	defer mockStore.Unlock()
	e, ok := mockStore.m[key]
	if !ok {
		return "", redis.Nil
	}
	// 检查过期
	if !e.expiry.IsZero() && time.Now().After(e.expiry) {
		delete(mockStore.m, key)
		return "", redis.Nil
	}
	// 删除并返回
	delete(mockStore.m, key)
	return e.val, nil
}
