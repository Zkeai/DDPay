package main

import (
	"context"
	"flag"
	cconf "github.com/Zkeai/DDPay/common/conf"
	"github.com/Zkeai/DDPay/common/cron"
	"github.com/Zkeai/DDPay/common/logger"
	"github.com/Zkeai/DDPay/common/utils"
	"github.com/Zkeai/DDPay/internal/conf"
	"github.com/Zkeai/DDPay/internal/watcher/okx"
	"github.com/Zkeai/DDPay/pkg/redis"
	"github.com/ouqiang/goutil"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

var (
	filePath *string
	AppDir   string
	LogDir   string
)

func main() {
	initEnv()

	// 启动 cron 服务
	cronService := cron.NewCronService()
	cronService.Start()

	// 启动初始化定时任务 获取汇率
	cronService.AddTask("*/30 * * * * *", func(ctx context.Context) {
		client := &okx.Redis{Redis: redis.GetClient()}
		client.GetOkxUsdtCnySellPrice()
		client.GetOkxTrxUsdtRate()

	})

	defer cronService.Stop()

	// 优雅退出处理
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	for {
		sig := <-ch
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			logger.Info("Shutting down CronService...")
			return
		default:
			return
		}
	}
}

func initEnv() {
	//logger 初始化
	logger.InitLogger()
	AppDir, err := goutil.WorkDir()
	if err != nil {
		logger.Fatal(err)
	}
	LogDir = filepath.Join(AppDir, "/log")
	utils.CreateDirIfNotExists(LogDir)

	// 读取配置
	filePath = flag.String("conf", "etc/config.yaml", "the config path")
	flag.Parse()

	c := new(conf.Conf)
	err = cconf.Unmarshal(*filePath, c)
	if err != nil {
		logger.Fatal("Failed to load config:", err)
	}
	//redis 初始化
	redis.InitRedis(c.Rdb)
	redis.GetClient().Set(context.Background(), "usdt-cny", "7.5", 0)
	redis.GetClient().Set(context.Background(), "usdt-trx", "0.27", 0)
}
