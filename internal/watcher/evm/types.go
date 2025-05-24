package evm

import (
	"github.com/go-redis/redis/v8"
	"time"
)

type TransferEvent struct {
	Chain     string
	TxHash    string
	From      string
	To        string
	Amount    string
	Timestamp time.Time
}

type Watcher struct {
	Chain          string
	RPC            string
	Contract       string
	Redis          *redis.Client
	KeyPrefix      string
	Callback       func(TransferEvent)
	PollDelay      time.Duration
	ConfirmBlocks  uint64
	TransferMethod string
}

type logEntry struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	TransactionHash  string   `json:"transactionHash"`
	BlockNumber      string   `json:"blockNumber"`
	TransactionIndex string   `json:"transactionIndex"`
	BlockHash        string   `json:"blockHash"`
	LogIndex         string   `json:"logIndex"`
	Removed          bool     `json:"removed"`
}
