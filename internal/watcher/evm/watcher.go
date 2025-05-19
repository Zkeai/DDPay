package evm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"
)

var transferEventTopic = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

func (w *Watcher) StartEvm() {
	go func() {
		ctx := context.Background()
		var fromBlock uint64 = 0

		for {
			latest, err := getLatestBlock(w.RPC)
			if err != nil {
				log.Printf("[%s] Error getting latest block: %v", w.Chain, err)
				time.Sleep(5 * time.Second)
				continue
			}

			toBlock := latest - 1
			if fromBlock == 0 {
				fromBlock = toBlock
			}

			keys, _ := w.Redis.Keys(ctx, w.KeyPrefix+"_*").Result()
			if len(keys) == 0 {
				time.Sleep(w.PollDelay)
				continue
			}

			watched := make(map[string]string) // address(lowercase) => redisKey
			for _, key := range keys {
				val, err := w.Redis.Get(ctx, key).Result()
				if err != nil {
					continue // 已过期
				}
				var entry struct {
					Address string `json:"address"`
					Status  string `json:"status"`
				}
				if err := json.Unmarshal([]byte(val), &entry); err != nil {
					continue
				}
				if entry.Status != "pending" {
					continue
				}
				watched[strings.ToLower(entry.Address)] = key
			}

			if len(watched) == 0 {
				time.Sleep(w.PollDelay)
				continue
			}

			logs, err := getLogs(w.RPC, "", fromBlock, toBlock)
			if err != nil {
				log.Printf("[%s] getLogs error: %v", w.Chain, err)
			} else {
				for _, logEntry := range logs {
					if len(logEntry.Topics) < 3 {
						continue
					}
					from := "0x" + logEntry.Topics[1][26:]
					to := "0x" + logEntry.Topics[2][26:]
					toLower := strings.ToLower(to)

					if redisKey, ok := watched[toLower]; ok {
						val := new(big.Int)
						val.SetString(strings.TrimPrefix(logEntry.Data, "0x"), 16)

						// 回调处理
						w.Callback(TransferEvent{
							Chain:     w.Chain,
							TxHash:    logEntry.TransactionHash,
							From:      from,
							To:        to,
							Amount:    val.String(),
							Timestamp: time.Now(),
						})

						// 更新 Redis 状态为 success
						newVal, _ := json.Marshal(map[string]string{
							"address": to,
							"status":  "success",
						})
						_ = w.Redis.Set(ctx, redisKey, newVal, time.Hour).Err()
					}
				}
			}

			fromBlock = toBlock + 1
			time.Sleep(w.PollDelay)
		}
	}()
}

func getLatestBlock(rpc string) (uint64, error) {
	req := map[string]interface{}{
		"jsonrpc": "2.0", "id": 1, "method": "eth_blockNumber", "params": []interface{}{},
	}
	body, _ := json.Marshal(req)
	resp, err := httpPost(rpc, body)
	if err != nil {
		return 0, err
	}
	var result struct{ Result string }
	if err := json.Unmarshal(resp, &result); err != nil {
		return 0, err
	}
	blockNum := new(big.Int)
	blockNum.SetString(strings.TrimPrefix(result.Result, "0x"), 16)
	return blockNum.Uint64(), nil
}

func getLogs(rpc string, contract string, from, to uint64) ([]logEntry, error) {
	params := map[string]interface{}{
		"fromBlock": fmt.Sprintf("0x%x", from),
		"toBlock":   fmt.Sprintf("0x%x", to),
		"topics":    [][]string{{transferEventTopic}},
	}
	if contract != "" {
		params["address"] = contract
	}
	req := map[string]interface{}{
		"jsonrpc": "2.0", "id": 1, "method": "eth_getLogs", "params": []interface{}{params},
	}
	body, _ := json.Marshal(req)
	resp, err := httpPost(rpc, body)
	if err != nil {
		return nil, err
	}
	var result struct{ Result []logEntry }
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return result.Result, nil
}

func httpPost(rpc string, body []byte) ([]byte, error) {
	resp, err := http.Post(rpc, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
