package service

import (
	"github.com/Zkeai/DDPay/internal/conf"
	"github.com/Zkeai/DDPay/pkg/telegram"
)
import "github.com/Zkeai/DDPay/internal/repo"

type Service struct {
	conf      *conf.Conf
	repo      *repo.Repo
	tgService *telegram.TelegramService
}

func NewService(conf *conf.Conf, tg *telegram.TelegramService) *Service {
	s := &Service{
		conf:      conf,
		repo:      repo.NewRepo(conf),
		tgService: tg,
	}

	return s

}
