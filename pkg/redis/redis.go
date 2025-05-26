package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"sync"
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
	once   sync.Once
)

func InitRedis(conf *Config) {
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     conf.Addr,
			Password: conf.Password, // No password set
			DB:       conf.Db,       // Use default DB
		})

		_, err := client.Ping(Ctx).Result()
		if err != nil {
			log.Fatalf("Failed to connect to Redis: %v", err)
		}
	})
}

func GetClient() *redis.Client {
	return client
}

// SetDefaultValues 设置默认值，仅在cron服务中调用
func SetDefaultValues() {
	if client == nil {
		log.Fatal("Redis client not initialized")
	}
	
	// 设置默认汇率
	client.Set(Ctx, "usdt-cny", "7.5", 0)
	client.Set(Ctx, "usdt-trx", "0.27", 0)
	
	log.Println("Redis default values set successfully")
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
