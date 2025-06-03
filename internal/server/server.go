package server

import (
	"log"

	chttp "github.com/Zkeai/DDPay/common/net/cttp"
	"github.com/Zkeai/DDPay/internal/conf"
	"github.com/Zkeai/DDPay/internal/handler"
	"github.com/Zkeai/DDPay/internal/service"
	"github.com/Zkeai/DDPay/pkg/email"
	"github.com/Zkeai/DDPay/pkg/telegram"
)

func NewHTTP(conf *conf.Conf, emailService *email.Service) *chttp.Server {
	// 创建HTTP服务器
	s := chttp.NewServer(conf.Server)
	
	// 设置Email服务（如果在Service中不能正常初始化）
	if conf.Email == nil && emailService != nil {
		// 这里可以手动设置conf.Email，如果需要的话
	}
	
	// 加载Telegram配置
	cfg, err := telegram.LoadConfig("etc/config.yaml")
	if err != nil {
		log.Fatalf("Telegram配置读取失败: %v", err)
	}

	// 初始化Telegram服务
	tgService, err := telegram.NewTelegramService(cfg)
	if err != nil {
		log.Fatalf("Telegram Bot初始化失败: %v", err)
	}
	
	// 设置Telegram服务到配置（如果在Service中不能正常初始化）
	conf.Tg = tgService.Bot
	
	// 注入service（使用新的参数签名）
	svc := service.NewService(conf)
	
	// 初始化路由
	handler.InitRouter(s, svc)

	// 启动Telegram监听处理
	go func() {
		tgService.StartHandler()
	}()

	// 启动HTTP服务
	err = s.Start()
	if err != nil {
		panic(err)
	}

	return s
}
