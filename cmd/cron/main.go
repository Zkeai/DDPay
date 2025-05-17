package main

import (
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	cconf "github.com/Zkeai/go_template/common/conf"
	"github.com/Zkeai/go_template/common/cron"
	"github.com/Zkeai/go_template/common/logger"
	"github.com/Zkeai/go_template/common/util"
	"github.com/Zkeai/go_template/internal/conf"
	"github.com/ouqiang/goutil"
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
	// logger 初始化
	logger.InitLogger()
	AppDir, _ = goutil.WorkDir()
	LogDir = filepath.Join(AppDir, "/log")
	util.CreateDirIfNotExists(LogDir)

	// 读取配置
	filePath = flag.String("conf", "etc/config.yaml", "the config path")
	flag.Parse()

	c := new(conf.Conf)
	err := cconf.Unmarshal(*filePath, c)
	if err != nil {
		logger.Fatal("Failed to load config:", err)
	}

}
