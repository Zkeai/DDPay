package main

import (
	"flag"
	"github.com/Zkeai/DDPay/internal/watcher"

	cconf "github.com/Zkeai/DDPay/common/conf"
	"github.com/Zkeai/DDPay/common/logger"
	"github.com/Zkeai/DDPay/common/utils"
	"github.com/Zkeai/DDPay/internal/conf"
	"github.com/Zkeai/DDPay/internal/wallet"
	yuredis "github.com/Zkeai/DDPay/pkg/redis"

	"github.com/ouqiang/goutil"
	"path/filepath"
)

var (
	filePath *string
	AppDir   string
	LogDir   string
)

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
	//初始化HD钱包
	if err := wallet.InitGlobalDeriver(c.Evm.Mnemonic); err != nil {
		logger.Error("初始化 HD 钱包失败", err)
	}

	//redis 初始化
	yuredis.InitRedis(c.Rdb)
	// 获取初始化后的 Redis 客户端
	client := yuredis.GetClient()

	//开启监控
	// 构造 RedisClient
	listener := watcher.NewRedisListener(client)

	// 启动监听
	go listener.WalletStart()

	select {} // 阻塞防止退出

}
