package main

import (
	"context"
	"flag"
	"github.com/Zkeai/DDPay/internal/server"
	"github.com/Zkeai/DDPay/pkg/redis"
	"os"
	"os/signal"
	"syscall"

	cconf "github.com/Zkeai/DDPay/common/conf"
	"github.com/Zkeai/DDPay/common/logger"
	"github.com/Zkeai/DDPay/common/utils"
	"github.com/Zkeai/DDPay/internal/conf"
	"github.com/ouqiang/goutil"
	"path/filepath"
)

var (
	// filePath yaml文件目录
	filePath *string
	// AppDir 应用根目录
	AppDir string
	// LogDir 日志目录
	LogDir string // 日志目录

)

// @title		DDPay API
// @version		1.0.0
// @description	DDpay https://github.com/zkeai/DDPay
// @host			localhost:2900
// @BasePath		/api/v1
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
	flag.Parse()
	filePath = flag.String("conf", "etc/config.yaml", "the config path")
	c := new(conf.Conf)
	err = cconf.Unmarshal(*filePath, c)
	if err != nil {
		logger.Error(err)
	}

	//redis 初始化
	redis.InitRedis(c.Rdb)

	//http 初始化
	srv := server.NewHTTP(c)
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
