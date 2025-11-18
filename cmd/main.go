package main

import (
	"time"

	"github.com/shy-robin/gochat/config"
	"github.com/shy-robin/gochat/pkg/global/log"
)

func main() {
	logConfig := config.GetConfig().Log
	log.InitLogger(logConfig.Path, logConfig.Level)

	for i := range 3 {
		if i == 0 {
			log.Logger.Sugar().Errorf("This is a test log at 0.")
		} else if i == 1 {
			log.Logger.Sugar().Debugf("This is a test log at 0.")
		} else {

			time.Sleep(time.Second)
			log.Logger.Sugar().Infof("This is a test log at %d.", i)
		}
	}
}
