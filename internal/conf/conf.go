package conf

import (
	"github.com/Zkeai/go_template/common/database"
	"github.com/Zkeai/go_template/common/mongodb"
	chttp "github.com/Zkeai/go_template/common/net/cttp"
	"github.com/Zkeai/go_template/common/redis"
	"github.com/Zkeai/go_template/pkg/solana/geysergrpc"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Conf struct {
	Server *chttp.Config          `yaml:"server"`
	DB     *database.Config       `yaml:"db"`
	Rdb    *redis.Config          `yaml:"redis"`
	Mongo  *mongodb.MongoConfig   `yaml:"mongo"`
	Tg     *tgbotapi.BotAPI       `yaml:"tg"`
	Grpc   *geysergrpc.GrpcConfig `yaml:"grpc"`
}
