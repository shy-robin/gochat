package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type TomlConfig struct {
	AppName string
	Log     LogConfig
	MySQL   MySQLConfig
	Api     ApiConfig
}

// 日志存储地址
type LogConfig struct {
	Path  string
	Level string
}

// 数据库配置
type MySQLConfig struct {
	Host        string
	Name        string
	Password    string
	Port        int
	TablePrefix string
	User        string
	Timeout     string
}

// 接口配置
type ApiConfig struct {
	Host   string
	Port   int
	Prefix string
}

var c TomlConfig

func InitConfig() {
	// 设置文件名
	viper.SetConfigName("config")
	// 设置文件类型
	viper.SetConfigType("toml")
	// 设置文件路径，可以多个viper会根据设置顺序依次查找
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	viper.Unmarshal(&c)
}

func GetConfig() TomlConfig {
	return c
}
