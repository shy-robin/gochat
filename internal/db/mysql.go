package db

import (
	"fmt"

	"github.com/shy-robin/gochat/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func init() {
	mysqlConfig := config.GetConfig().MySQL
	username := mysqlConfig.User     // 账号
	password := mysqlConfig.Password // 密码
	host := mysqlConfig.Host         // 数据库地址，可以是 Ip 或者域名
	port := mysqlConfig.Port         // 数据库端口
	dbName := mysqlConfig.Name       // 数据库名
	timeout := mysqlConfig.Timeout   // 连接超时，10s

	// 拼接下 dsn 参数, dsn 格式可以参考上面的语法，这里使用 Sprintf 动态拼接 dsn 参数，因为一般数据库连接参数，我们都是保存在配置文件里面，需要从配置文件加载参数，然后拼接dsn。
	// 构造一个不带数据库名的 DSN，用于连接到 MySQL 服务器本身
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8&parseTime=True&loc=Local&timeout=%s", username, password, host, port, timeout)

	// 第一次连接：连接到 MySQL 服务器，而不是特定的数据库
	tempDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(fmt.Errorf("无法连接到 MySQL 服务器: %w", err))
	}

	// 执行 SQL 语句创建数据库
	createDbSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", dbName)
	if result := tempDB.Exec(createDbSql); result.Error != nil {
		panic(fmt.Errorf("创建数据库 %s 失败: %w", dbName, result.Error))
	}

	// 构造带数据库名的完整 DSN
	appDsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", username, password, host, port, dbName, timeout)
	db, err = gorm.Open(mysql.Open(appDsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(fmt.Errorf("连接应用数据库 %s 失败: %w", dbName, err))
	}

	sqlDB, _ := db.DB()

	// 设置数据库连接池参数
	sqlDB.SetMaxOpenConns(100) // 设置数据库连接池最大连接数
	sqlDB.SetMaxIdleConns(20)  // 连接池最大允许的空闲连接数，如果没有 sql 任务需要执行的连接数大于 20，超过的连接会被连接池关闭。
}

func GetDB() *gorm.DB {
	return db
}
