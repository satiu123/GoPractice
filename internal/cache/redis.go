package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var ctx = context.Background()

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// InitRedis 初始化Redis连接
func InitRedis(config RedisConfig) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	// 测试连接
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %v", err)
	}

	log.Println("Redis connected successfully")
	return nil
}

// SetCache 设置缓存
func SetCache(key string, value any, expiration time.Duration) error {
	// 将值序列化为JSON
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %v", err)
	}

	return RedisClient.Set(ctx, key, jsonValue, expiration).Err()
}

// GetCache 获取缓存
func GetCache(key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}

// GetCacheJSON 获取缓存并反序列化为指定类型
func GetCacheJSON(key string, dest any) error {
	val, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// DeleteCache 删除缓存
func DeleteCache(key string) error {
	return RedisClient.Del(ctx, key).Err()
}

// Exists 检查键是否存在
func Exists(key string) (bool, error) {
	result, err := RedisClient.Exists(ctx, key).Result()
	return result > 0, err
}

// SetCacheWithContext 使用上下文设置缓存
func SetCacheWithContext(ctx context.Context, key string, value any, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %v", err)
	}

	return RedisClient.Set(ctx, key, jsonValue, expiration).Err()
}

// GetCacheWithContext 使用上下文获取缓存
func GetCacheWithContext(ctx context.Context, key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}
