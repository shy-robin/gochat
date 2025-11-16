package main

import (
	"time"

	"github.com/shy-robin/gochat/config"
	"github.com/shy-robin/gochat/pkg/global/log"
)

func main() {
	logConfig := config.GetConfig().Log
	log.InitLogger(logConfig.Path, logConfig.Level)

	for i := range 10 {
		time.Sleep(time.Second)
		log.Logger.Sugar().Infof("This is a test log at %d.", i)
	}
}
