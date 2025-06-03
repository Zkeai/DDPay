package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/Zkeai/DDPay/internal/server"
	"github.com/Zkeai/DDPay/internal/wallet"
	"github.com/Zkeai/DDPay/pkg/email"
	"github.com/Zkeai/DDPay/pkg/redis"

	"path/filepath"

	cconf "github.com/Zkeai/DDPay/common/conf"
	"github.com/Zkeai/DDPay/common/logger"
	"github.com/Zkeai/DDPay/common/utils"
	"github.com/Zkeai/DDPay/internal/conf"
	"github.com/ouqiang/goutil"
)

var (
	// filePath yaml文件目录
	filePath *string
	// AppDir 应用根目录
	AppDir string
	// LogDir 日志目录
	LogDir string // 日志目录

)

// @title	DDPay API
// @version	1.0.0
// @description	DDpay https://github.com/zkeai/DDPay
// @host	localhost:2900
// @BasePath	/api/v1
func main() {
	//logger 初始化
	logger.InitLogger()
	AppDir, err := goutil.WorkDir()
	if err != nil {
		logger.Fatal(err)
	}
	LogDir = filepath.Join(AppDir, "/log")
	utils.CreateDirIfNotExists(LogDir)

	//读取yaml配置
	filePath = flag.String("conf", "etc/config.yaml", "the config path")
	flag.Parse()
	c := new(conf.Conf)
	err = cconf.Unmarshal(*filePath, c)
	if err != nil {
		logger.Error(err)
	}

	//全局化初始化
	conf.Load(c)
	//redis 初始化
	redis.InitRedis(c.Rdb)

	//wallet 初始化
	//初始化HD钱包
	if err := wallet.InitGlobalDeriver(c.Evm.Mnemonic); err != nil {
		logger.Error("初始化 HD 钱包失败", err)
	}
	
	//邮件服务初始化
	emailService := email.NewService(*c.Email)
	logger.Info("初始化邮件服务成功")

	//http 初始化
	srv := server.NewHTTP(c, emailService)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-ch
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			_ = srv.Shutdown(context.Background())
			return
		default:

			return
		}
	}
}
