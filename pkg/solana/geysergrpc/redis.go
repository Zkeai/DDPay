package geysergrpc

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisListener struct {
	client *redis.Client
}

func NewRedisListener(rdb *redis.Client) *RedisListener {
	return &RedisListener{client: rdb}
}

// Start 启动监听 Redis 中的订阅配置变化
func (r *RedisListener) grpcStart() {
	pubsub := r.client.PSubscribe(ctx, "__keyspace@0__:user_grpc_*") // 使用 keyspace 监听
	defer func(pubsub *redis.PubSub) {
		err := pubsub.Close()
		if err != nil {

		}
	}(pubsub)

	log.Println("[RedisListener] 开始监听用户订阅配置...")

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			log.Println("[RedisListener] 接收消息失败:", err)
			continue
		}

		// 只处理 set 操作
		if strings.HasPrefix(msg.Channel, "__keyspace@0__:user_grpc_") && msg.Payload == "set" {
			userID := strings.TrimPrefix(msg.Channel, "__keyspace@0__:user_grpc_")
			log.Printf("✅ [RedisListener] 检测到用户 %s 的订阅配置发生变化 (key: %s)\n", userID, "user_grpc_"+userID)
			go r.handleUserConfig(userID)
		}
	}
}

// handleUserConfig 处理订阅配置变更
func (r *RedisListener) handleUserConfig(userID string) {
	key := "user_grpc_" + userID

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		log.Printf("[RedisListener] 获取用户 %s 配置失败: %v\n", userID, err)
		return
	}

	var cfg SubscribeConfig
	if err := json.Unmarshal([]byte(val), &cfg); err != nil {
		log.Printf("[RedisListener] 解析用户 %s 配置失败: %v\n", userID, err)
		return
	}

	go NewGeyserClient(cfg)

	log.Printf("[Geyser] 启动订阅任务：用户 %s\n", userID)

}
