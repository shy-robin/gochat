package log

import (
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// zap 的使用参考：bilibili.com/video/BV1Rk99YHEM6

const (
	BLUE_COLOR   = "\033[0;34m"
	YELLOW_COLOR = "\033[0;33m"
	RED_COLOR    = "\033[0;31m"
	RESET_COLOR  = "\033[0m"
)

var (
	Logger *zap.Logger
	// 将配置对象的完整内容，以键值对的形式，嵌入到日志记录中
	Any    = zap.Any
	String = zap.String
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
	upperLevel := strings.ToUpper(level.String())

	switch level {
	case zapcore.InfoLevel:
		enc.AppendString(BLUE_COLOR + upperLevel + RESET_COLOR)
	case zapcore.WarnLevel:
		enc.AppendString(YELLOW_COLOR + upperLevel + RESET_COLOR)
	case zapcore.ErrorLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		enc.AppendString(RED_COLOR + upperLevel + RESET_COLOR)
	default:
		enc.AppendString(upperLevel)
	}
}

// logEncoder 时间分片和 level 分片同时做
type logEncoder struct {
	zapcore.Encoder
	logPath     string
	logLevel    string
	file        *os.File
	errFile     *os.File
	currentDate string
}

func (this *logEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	// 先调用原始的 EncodeEntry 方法生成日志行
	buff, err := this.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return nil, err
	}

	data := buff.String()
	buff.Reset()
	buff.AppendString("[myApp]" + data)
	data = buff.String()

	// 时间分片
	now := time.Now().Format("2006-01-02")
	dirName := fmt.Sprintf("%s/%s", this.logPath, now)

	if this.currentDate != now {
		mkdirErr := os.MkdirAll(dirName, 0755)
		if mkdirErr != nil {
			return nil, mkdirErr
		}
		fileName := fmt.Sprintf("%s/output.log", dirName)
		file, _ := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
		this.file = file
		this.currentDate = now
	}

	switch entry.Level {
	case zapcore.ErrorLevel:
		if this.errFile == nil {
			fileName := fmt.Sprintf("%s/error.log", dirName)
			file, _ := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
			this.errFile = file
		}
		this.errFile.WriteString(buff.String())
	}

	if this.currentDate == now {
		this.file.WriteString(data)
	}

	return buff, nil
}

func InitLogger(logPath string, logLevel string) *zap.Logger {
	// 测试代码
	// devLog()
	// exampleLog()
	// prodLog()

	// logger, _ := zap.NewDevelopment()

	// 使用 zap 的 NewDevelopmentConfig 快速配置
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

	// 根据传入的日志级别设置配置
	var level zapcore.Level
	switch logLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "panic":
		level = zapcore.PanicLevel
	case "fatal":
		level = zapcore.FatalLevel
	default:
		level = zapcore.InfoLevel // 默认为 info 级别
	}
	cfg.Level = zap.NewAtomicLevelAt(level)

	// 创建自定义的 Encoder
	encoder := &logEncoder{
		Encoder:  zapcore.NewConsoleEncoder(cfg.EncoderConfig), // 使用 Console 编码器
		logPath:  logPath,
		logLevel: logLevel,
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout), // 输出到控制台
		level,                      // 设置显示的日志级别（debug 可以打印出 debug info warn）
	)

	// 创建 logger 实例
	logger := zap.New(core, zap.AddCaller()) // 显示堆栈跟踪信息

	// 全局替换 zap 实例，使 zap.L().Info 能够调用
	zap.ReplaceGlobals(logger)

	Logger = logger

	return logger
}
