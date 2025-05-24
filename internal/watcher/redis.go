package watcher

import (
	"context"
	"encoding/json"
	"github.com/Zkeai/DDPay/internal/watcher/evm"
	"log"
	"strconv"
	"strings"
	"time"

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

	var info conf.SubscribeConfig
	if err := json.Unmarshal([]byte(val), &info); err != nil {
		log.Printf("[RedisWalletListener] 解析配置失败 (%s): %v\n", key, err)
		return
	}

	//启动wallet监控
	for _, cfg := range conf.EVMChains {
		if cfg.Name == "bsc" {
			bscWatcher := &evm.Watcher{
				Chain:     cfg.Name,
				RPC:       cfg.RPC,
				Contract:  cfg.ContractAddr,
				Redis:     r.client,
				KeyPrefix: cfg.Name + "_" + strconv.FormatInt(info.MerchantID, 10),
				Callback: func(e evm.TransferEvent) {
					log.Printf("[CALLBACK][%s] Tx: %s, From: %s, To: %s, Amount: %s\n", e.Chain, e.TxHash, e.From, e.To, e.Amount)
					// 你可以在这里触发数据库写入、Webhook 回调等逻辑
					
				},
				PollDelay:      time.Second * 2,
				ConfirmBlocks:  20,
				TransferMethod: "0xa9059cbb", // transfer(address,uint256)
			}
			bscWatcher.StartEvm()
		}

	}

}
