package server

import (
	chttp "github.com/Zkeai/DDPay/common/net/cttp"
	"github.com/Zkeai/DDPay/internal/conf"
	"github.com/Zkeai/DDPay/internal/handler"
	"github.com/Zkeai/DDPay/internal/service"
	"github.com/Zkeai/DDPay/pkg/telegram"
	"log"
)

func NewHTTP(conf *conf.Conf) *chttp.Server {

	s := chttp.NewServer(conf.Server)
	//telegram
	cfg, err := telegram.LoadConfig("etc/config.yaml")
	if err != nil {
		log.Fatalf("配置读取失败: %v", err)
	}

	tgService, err := telegram.NewTelegramService(cfg)
	if err != nil {
		log.Fatalf("Bot 初始化失败: %v", err)
	}

	// 注入 service
	svc := service.NewService(conf, tgService)
	handler.InitRouter(s, svc)

	// 启动监听处理

	go func() {
		tgService.StartHandler()
	}()

	err = s.Start()

	if err != nil {
		panic(err)
	}

	return s
}
