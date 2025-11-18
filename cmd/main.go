package main

import (
	"github.com/shy-robin/gochat/config"
	"github.com/shy-robin/gochat/pkg/global/log"
)

func main() {
	// 初始化 Logger
	logConfig := config.GetConfig().Log
	log.InitLogger(logConfig.Path, logConfig.Level)
	log.Logger.Info("config", log.Any("config", config.GetConfig()))
}
