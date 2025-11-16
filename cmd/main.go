package main

import (
	"fmt"
	"github.com/shy-robin/gochat/config"
)

func main() {
	fmt.Println("Hello, gochat!")
	log := config.GetConfig().Log

	fmt.Println(log)
}
