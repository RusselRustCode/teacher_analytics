package redis

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/redis/go-redis/v9"
    
    "github.com/RusselRustCode/teacher_analytics/core-service/internal/interfaces"
)

type RedisCache struct {
    client *redis.Client
}

func NewRedisCache(addr, password string, db int) interfaces.Cache {
    client := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })
    
    return &RedisCache{
        client: client,
    }
}

func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
    val, err := c.client.Get(ctx, key).Result()
    if err == redis.Nil {
        return "", nil
    }
    return val, err
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
    var data string
    
    switch v := value.(type) {
    case string:
        data = v
    case []byte:
        data = string(v)
    default:
        jsonData, err := json.Marshal(v)
        if err != nil {
            return fmt.Errorf("failed to marshal value: %w", err)
        }
        data = string(jsonData)
    }
    
    return c.client.Set(ctx, key, data, expiration).Err()
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
    return c.client.Del(ctx, key).Err()
}

func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
    count, err := c.client.Exists(ctx, key).Result()
    return count > 0, err
}

func (c *RedisCache) Close() error {
    return c.client.Close()
}