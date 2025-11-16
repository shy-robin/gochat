package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zap 的使用参考：bilibili.com/video/BV1Rk99YHEM6

const (
	BLUE   = "\033[0;34m"
	YELLOW = "\033[0;33m"
	RED    = "\033[0;31m"
	RESET  = "\033[0m"
)

var (
	Logger *zap.Logger
)

func devLog() {
	logger, _ := zap.NewDevelopment()

	logger.Debug("This is DEV Debug log.")
	logger.Info("This is DEV Info log.")
	// 以下方法会打印调用栈
	logger.Warn("This is DEV Warn log.")
	logger.Error("This is DEV Error log.")
	// 以下方法会抛出错误，并不再执行后续代码
	logger.Panic("This is DEV Panic log.")
	logger.Fatal("This is DEV Fatal log.")
}

func exampleLog() {
	// 日志输出使用 json 格式
	logger := zap.NewExample()

	logger.Debug("This is Example Debug log.")
	logger.Info("This is Example Info log.")
	// 以下方法会打印调用栈
	logger.Warn("This is Example Warn log.")
	logger.Error("This is Example Error log.")
	// 以下方法会抛出错误，并不再执行后续代码
	logger.Panic("This is Example Panic log.")
	logger.Fatal("This is Example Fatal log.")
}

func prodLog() {
	// 日志输出使用 json 格式
	logger, _ := zap.NewProduction()

	logger.Debug("This is Prod Debug log.")
	logger.Info("This is Prod Info log.")
	// 以下方法会打印调用栈
	logger.Warn("This is Prod Warn log.")
	logger.Error("This is Prod Error log.")
	// 以下方法会抛出错误，并不再执行后续代码
	logger.Panic("This is Prod Panic log.")
	logger.Fatal("This is Prod Fatal log.")
}

func levelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch level {
	case zapcore.InfoLevel:
		enc.AppendString(BLUE + level.String() + RESET)
	case zapcore.WarnLevel:
		enc.AppendString(YELLOW + level.String() + RESET)
	case zapcore.ErrorLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		enc.AppendString(RED + level.String() + RESET)
	default:
		enc.AppendString(level.String())
	}
}

func InitLogger(logPath string, logLevel string) {
	// 测试代码
	// devLog()
	// exampleLog()
	// prodLog()

	// logger, _ := zap.NewDevelopment()

	cfg := zap.NewDevelopmentConfig()

	// debug 可以打印出 debug info warn
	// info  级别可以打印 info warn
	// warn  只能打印 warn
	// debug->info->warn->error
	// cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	// 格式化时间
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	// 输出美化颜色
	cfg.EncoderConfig.EncodeLevel = levelEncoder

	// 控制台日志
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg.EncoderConfig),
		zapcore.AddSync(os.Stdout), // 输出到控制台
		zapcore.DebugLevel,         // 设置显示的日志级别（debug 可以打印出 debug info warn）
	)
	// 文件日志
	file, _ := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	fileCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg.EncoderConfig),
		zapcore.AddSync(file), // 输出到文件
		zapcore.DebugLevel,    // 设置显示的日志级别（debug 可以打印出 debug info warn）
	)

	// 日志双写（控制台 + 文件）
	core := zapcore.NewTee(consoleCore, fileCore)

	// 创建 logger 实例
	logger := zap.New(core, zap.AddCaller()) // 显示堆栈跟踪信息

	Logger = logger

	// 全局替换 zap 实例，使 zap.L().Info 能够调用
	zap.ReplaceGlobals(logger)
}
