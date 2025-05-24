package evm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Zkeai/DDPay/common/logger"
	"io"
	"math/big"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/smallnest/chanx"
)

var (
	lastBlockNumbers = make(map[string]uint64)
	blockScanQueues  = make(map[string]*chanx.UnboundedChan[uint64])
	scanTotals       = make(map[string]*uint64)
	scanSuccesses    = make(map[string]*uint64)
)

func (w *Watcher) StartEvm() {
	// 初始化链相关的变量
	if _, ok := lastBlockNumbers[w.Chain]; !ok {
		lastBlockNumbers[w.Chain] = 0
		blockScanQueues[w.Chain] = chanx.NewUnboundedChan[uint64](context.Background(), 30)
		var total, success uint64
		scanTotals[w.Chain] = &total
		scanSuccesses[w.Chain] = &success

	}

	// 启动区块扫描工作池
	go w.startBlockScan()
	// 启动区块高度监控
	go w.monitorBlockHeight()

}

func (w *Watcher) startBlockScan() {
	p, err := ants.NewPoolWithFunc(8, w.processBlock)
	if err != nil {
		logger.Error("[%s] 创建工作池失败: %v", w.Chain, err)
		return
	}
	defer p.Release()

	queue := blockScanQueues[w.Chain]
	for blockNum := range queue.Out {
		if err := p.Invoke(blockNum); err != nil {
			queue.In <- blockNum
			logger.Error("[%s] 处理区块失败: %v", w.Chain, err)
		}
	}
}

func (w *Watcher) processBlock(n interface{}) {
	blockNum := n.(uint64)
	atomic.AddUint64(scanTotals[w.Chain], 1)

	// 获取区块信息
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params":  []interface{}{fmt.Sprintf("0x%x", blockNum), true},
		"id":      1,
	}

	body, _ := json.Marshal(req)
	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Post(w.RPC, "application/json", bytes.NewReader(body))
	if err != nil {
		blockScanQueues[w.Chain].In <- blockNum
		logger.Error("[%s] 获取区块失败: %d, 错误: %v", w.Chain, blockNum, err)
		return
	}

	blockData, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		blockScanQueues[w.Chain].In <- blockNum
		logger.Error("[%s] 读取区块数据失败: %d, 错误: %v", w.Chain, blockNum, err)
		return
	}

	var result struct {
		Result struct {
			Timestamp    string `json:"timestamp"`
			Transactions []struct {
				To    string `json:"to"`
				Input string `json:"input"`
				From  string `json:"from"`
				Hash  string `json:"hash"`
			} `json:"transactions"`
		} `json:"result"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(blockData, &result); err != nil {
		blockScanQueues[w.Chain].In <- blockNum
		logger.Error("[%s] 解析区块数据失败: %d, 错误: %v", w.Chain, blockNum, err)
		return
	}

	if result.Error != nil {
		blockScanQueues[w.Chain].In <- blockNum
		logger.Error("[%s] 区块错误: %d, 错误: %s", w.Chain, blockNum, result.Error.Message)
		return
	}

	// 处理交易
	ctx := context.Background()
	keys, _ := w.Redis.Keys(ctx, w.KeyPrefix+"_*").Result()
	watched := make(map[string]string)

	for _, key := range keys {
		val, err := w.Redis.Get(ctx, key).Result()
		if err != nil {
			logger.Error("[%s] 获取Redis键失败: %s, 错误: %v", w.Chain, key, err)
			continue
		}

		var entry struct {
			Address string `json:"address"`
			Status  string `json:"status"`
		}

		if err := json.Unmarshal([]byte(val), &entry); err != nil {
			logger.Error("[%s] 解析Redis数据失败: %s, 错误: %v", w.Chain, key, err)
			continue
		}

		if entry.Status == "pending" {
			watched[strings.ToLower(entry.Address)] = key

		}
	}

	timestamp, _ := new(big.Int).SetString(strings.TrimPrefix(result.Result.Timestamp, "0x"), 16)

	for _, tx := range result.Result.Transactions {

		if tx.To != strings.ToLower(w.Contract) {

			continue
		}

		if len(tx.Input) < 10 {

			continue
		}

		// 检查是否是转账方法
		if !strings.HasPrefix(tx.Input, w.TransferMethod) {
			continue
		}

		// 解析转账数据
		from := tx.From
		to := "0x" + tx.Input[34:74]
		amount := new(big.Int)
		amount.SetString(tx.Input[74:], 16)

		toLower := strings.ToLower(to)
		if redisKey, ok := watched[toLower]; ok {

			w.Callback(TransferEvent{
				Chain:     w.Chain,
				TxHash:    tx.Hash,
				From:      from,
				To:        to,
				Amount:    amount.String(),
				Timestamp: time.Unix(timestamp.Int64(), 0),
			})

			// 更新 Redis 状态
			newVal, _ := json.Marshal(map[string]string{
				"address": to,
				"status":  "success",
			})
			err := w.Redis.Set(ctx, redisKey, newVal, time.Hour).Err()
			if err != nil {
				return
			}

		}
	}

	atomic.AddUint64(scanSuccesses[w.Chain], 1)

}

func (w *Watcher) monitorBlockHeight() {

	for range time.Tick(w.PollDelay) {
		latest, err := getLatestBlock(w.RPC)
		if err != nil {

			continue
		}

		// 考虑确认区块数
		if w.ConfirmBlocks > 0 {
			latest = latest - w.ConfirmBlocks
		}

		// 首次启动
		if lastBlockNumbers[w.Chain] == 0 {
			lastBlockNumbers[w.Chain] = latest - 1

		}

		// 区块高度没有变化
		if latest <= lastBlockNumbers[w.Chain] {
			continue
		}

		// 将新区块加入队列
		queue := blockScanQueues[w.Chain]
		for n := lastBlockNumbers[w.Chain] + 1; n <= latest; n++ {
			queue.In <- n
		}

		lastBlockNumbers[w.Chain] = latest
	}
}

func getLatestBlock(rpc string) (uint64, error) {
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_blockNumber",
		"params":  []interface{}{},
		"id":      1,
	}

	body, _ := json.Marshal(req)
	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Post(rpc, "application/json", bytes.NewReader(body))
	if err != nil {
		return 0, err
	}

	blockData, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return 0, err
	}

	var result struct {
		Result string `json:"result"`
		Error  *struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(blockData, &result); err != nil {
		return 0, err
	}

	if result.Error != nil {
		return 0, fmt.Errorf("RPC error: %s", result.Error.Message)
	}

	blockNum := new(big.Int)
	blockNum.SetString(strings.TrimPrefix(result.Result, "0x"), 16)
	return blockNum.Uint64(), nil
}
