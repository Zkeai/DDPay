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
	Chain     string
	RPC       string
	Redis     *redis.Client
	PollDelay time.Duration
	Callback  func(evt TransferEvent)
	KeyPrefix string // like "bsc" or "arb"
}

type logEntry struct {
	Address         string   `json:"address"`
	Topics          []string `json:"topics"`
	Data            string   `json:"data"`
	TransactionHash string   `json:"transactionHash"`
}
