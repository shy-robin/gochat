package main

import (
	"fmt"

	"github.com/shy-robin/gochat/config"
	"github.com/shy-robin/gochat/internal/db"
	"github.com/shy-robin/gochat/internal/router"
	"github.com/shy-robin/gochat/pkg/common"
	"github.com/shy-robin/gochat/pkg/global/log"
)

func main() {
	// 初始化配置文件
	config.InitConfig()

	// 初始化 Logger
	logConfig := config.GetConfig().Log
	log.InitLogger(logConfig.Path, logConfig.Level)
	log.Logger.Info("config", log.Any("config", config.GetConfig()))

	// 初始化数据库
	db.InitMysqlDB()

	// 初始化路由
	ginServer := router.NewRouter()
	common.SetupCustomValidator(ginServer)

	// 启动服务
	apiConfig := config.GetConfig().Api
	ginServer.Run(fmt.Sprintf("%s:%d", apiConfig.Host, apiConfig.Port))
}
