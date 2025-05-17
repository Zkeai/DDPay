package server

import (
	"github.com/Zkeai/go_template/internal/handler"
	"github.com/Zkeai/go_template/internal/service"
	ws2 "github.com/Zkeai/go_template/internal/ws"
	"github.com/Zkeai/go_template/pkg/telegram"
	"log"
	"time"

	chttp "github.com/Zkeai/go_template/common/net/cttp"
	"github.com/Zkeai/go_template/internal/conf"
	"github.com/Zkeai/go_template/internal/wsHandler"
)

func NewHTTP(conf *conf.Conf) *chttp.Server {
	// 注册模块处理函数
	ws2.RegisterModuleHandler("coin", wsHandler.CoinHandler)
	ws2.RegisterModuleHandler("chat", wsHandler.ChatHandler)

	s := chttp.NewServer(conf.Server)
	hub := ws2.NewHub()
	go hub.Run()
	hub.StartHealthChecker(2 * time.Minute)

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
	handler.InitRouter(s, svc, hub)
	telegram.RegisterUpsertChannelHandler(svc.UpsertChannel)

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
