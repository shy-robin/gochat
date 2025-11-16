package main

import (
	"github.com/shy-robin/gochat/config"
	"github.com/shy-robin/gochat/pkg/global/log"
)

func main() {
	logConfig := config.GetConfig().Log
	log.InitLogger(logConfig.Path, logConfig.Level)
}
