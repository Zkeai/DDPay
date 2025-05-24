package conf

import (
	"github.com/Zkeai/DDPay/common/database"
	chttp "github.com/Zkeai/DDPay/common/net/cttp"
	"github.com/Zkeai/DDPay/internal/wallet"
	"github.com/Zkeai/DDPay/pkg/redis"
	"github.com/Zkeai/DDPay/pkg/solana/geysergrpc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Conf struct {
	DB     *database.Config       `yaml:"db"`
	Server *chttp.Config          `yaml:"server"`
	Rdb    *redis.Config          `yaml:"redis"`
	Tg     *tgbotapi.BotAPI       `yaml:"tg"`
	Grpc   *geysergrpc.GrpcConfig `yaml:"grpc"`
	Evm    *wallet.Config         `yaml:"evm"`
}

type ChainConfig struct {
	Name         string
	RPC          string
	ContractAddr string
}

var EVMChains = []ChainConfig{
	//{
	//	Name:         "bsc",
	//	RPC:          "https://bsc-dataseed.binance.org",
	//	ContractAddr: "0x55d398326f99059fF775485246999027B3197955",
	//},
	{
		Name:         "pol",
		RPC:          "https://polygon-rpc.com",
		ContractAddr: "0xc2132d05d31c914a87c6611c10748aeb04b58e8f",
	},
	{
		Name:         "arb",
		RPC:          "https://arb1.arbitrum.io/rpc",
		ContractAddr: "0xfd086bc7cd5c481dcc9c85ebe478a1c0b69fcbb9",
	},
	{
		Name:         "bsc",
		RPC:          "https://data-seed-prebsc-1-s3.binance.org:8545",
		ContractAddr: "0x337610d27c682E347C9cD60BD4b3b107C9d34dDd",
	},
}

type SubscribeConfig struct {
	OrderID    string
	Chain      string
	MerchantID int64
	Amount     float64
	Wallet     string
	Status     string
	ExpireTime int
}
