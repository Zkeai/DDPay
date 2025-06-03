package conf

import (
	"github.com/Zkeai/DDPay/common/conf"
	"github.com/Zkeai/DDPay/common/database"
	chttp "github.com/Zkeai/DDPay/common/net/cttp"
	"github.com/Zkeai/DDPay/internal/wallet"
	"github.com/Zkeai/DDPay/pkg/email"
	"github.com/Zkeai/DDPay/pkg/jwt"
	"github.com/Zkeai/DDPay/pkg/redis"
	"github.com/Zkeai/DDPay/pkg/solana/geysergrpc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// OAuthConfig GitHub和Google OAuth配置结构体
type OAuthConfig struct {
	Github GithubConfig `yaml:"github"`
	Google GoogleConfig `yaml:"google"`
}

// GithubConfig GitHub OAuth配置
type GithubConfig struct {
	ClientID     string `yaml:"clientID"`
	ClientSecret string `yaml:"clientSecret"`
	RedirectURI  string `yaml:"redirectURI"`
	Scopes       string `yaml:"scopes"`
}

// GoogleConfig Google OAuth配置
type GoogleConfig struct {
	ClientID     string `yaml:"clientID"`
	ClientSecret string `yaml:"clientSecret"`
	RedirectURI  string `yaml:"redirectURI"`
	Scopes       string `yaml:"scopes"`
}

type Conf struct {
	DB     *database.Config       `yaml:"db"`
	Server *chttp.Config          `yaml:"server"`
	Rdb    *redis.Config          `yaml:"redis"`
	Tg     *tgbotapi.BotAPI       `yaml:"tg"`
	Grpc   *geysergrpc.GrpcConfig `yaml:"grpc"`
	Evm    *wallet.Config         `yaml:"evm"`
	Config *conf.Config           `yaml:"config"`
	JWT    *jwt.Config            `yaml:"jwt"`
	Email  *email.Config          `yaml:"email"`
	OAuth  *OAuthConfig           `yaml:"oauth"` // 添加OAuth配置
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
var SignKey = "DDPay"

type SubscribeConfig struct {
	OrderID    string
	Chain      string
	MerchantID int64
	Amount     float64
	Wallet     string
	Status     string
	ExpireTime int
}

var globalConf *Conf

func Load(c *Conf) {
	globalConf = c
}

func Get() *Conf {
	return globalConf
}
