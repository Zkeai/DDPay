package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type Config struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}

var (
	client *redis.Client
	Ctx    = context.Background()
)

func InitRedis(conf *Config) {
	client = redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password, // No password set
		DB:       conf.Db,       // Use default DB
	})

	_, err := client.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}

func GetClient() *redis.Client {
	return client
}

func Set(key string, value interface{}, expiration time.Duration) error {
	return client.Set(Ctx, key, value, expiration).Err()
}

func Get(key string) (string, error) {
	return client.Get(Ctx, key).Result()
}

func Del(keys ...string) error {
	return client.Del(Ctx, keys...).Err()
}

func Exists(keys ...string) (int64, error) {
	return client.Exists(Ctx, keys...).Result()
}

// ScanKeys 模糊匹配 key（用于查找所有订单）
func ScanKeys(pattern string) ([]string, error) {
	var (
		cursor uint64
		keys   []string
		err    error
	)

	allKeys := make([]string, 0)
	for {
		keys, cursor, err = client.Scan(Ctx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, err
		}
		allKeys = append(allKeys, keys...)
		if cursor == 0 {
			break
		}
	}
	return allKeys, nil
}
