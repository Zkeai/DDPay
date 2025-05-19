package watcher

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/Zkeai/DDPay/internal/conf"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisListener struct {
	client *redis.Client
}

func NewRedisListener(rdb *redis.Client) *RedisListener {
	return &RedisListener{client: rdb}
}

// WalletStart Start 启动监听 Redis 中的订阅配置变化
func (r *RedisListener) WalletStart() {
	pubsub := r.client.PSubscribe(ctx,
		"__keyspace@0__:bsc_*_*",
		"__keyspace@0__:arb_*_*",
		"__keyspace@0__:sol_*_*",
		"__keyspace@0__:pol_*_*",
	)
	defer func(pubsub *redis.PubSub) {
		if err := pubsub.Close(); err != nil {
			log.Printf("[RedisWalletListener] 关闭订阅失败: %v\n", err)
		}
	}(pubsub)

	log.Println("[RedisWalletListener] 开始监听多链用户订阅配置...")

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			log.Println("[RedisWalletListener] 接收消息失败:", err)
			continue
		}

		if msg.Payload != "set" {
			continue
		}

		keyParts := strings.SplitN(msg.Channel, "__keyspace@0__:", 2)
		if len(keyParts) < 2 {
			continue
		}

		key := keyParts[1]
		log.Printf("✅ [RedisWalletListener] 检测到订阅配置变化: %s\n", key)
		go r.handleChainConfig(key)
	}
}

// handleWalletConfig 处理订阅配置变更
func (r *RedisListener) handleChainConfig(key string) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		log.Printf("[RedisWalletListener] 获取配置失败 (%s): %v\n", key, err)
		return
	}

	var cfg conf.SubscribeConfig
	if err := json.Unmarshal([]byte(val), &cfg); err != nil {
		log.Printf("[RedisWalletListener] 解析配置失败 (%s): %v\n", key, err)
		return
	}

	//启动wallet监控

}
