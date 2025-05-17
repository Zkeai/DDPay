package telegram

import (
	"fmt"

	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v3"
)

type Config struct {
	BotToken string `yaml:"botToken"`
}

type TelegramService struct {
	Bot *tgbotapi.BotAPI
}

// LoadConfig 从 config.yaml 中加载 Telegram 配置
func LoadConfig(path string) (*Config, error) {
	cfg := struct {
		Telegram Config `yaml:"telegram"`
	}{}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg.Telegram, nil
}

// NewTelegramService 初始化 Telegram Bot 实例
func NewTelegramService(cfg *Config) (*TelegramService, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)
	return &TelegramService{Bot: bot}, nil
}
