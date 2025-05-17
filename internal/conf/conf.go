package conf

import (
	"github.com/Zkeai/DDPay/common/database"
	chttp "github.com/Zkeai/DDPay/common/net/cttp"
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
}
