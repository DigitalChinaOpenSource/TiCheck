package logutil

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

const (
	// DefaultLogFilePath is the default path for saving log files.
	DefaultLogFilePath = "../../log/ticheck.log"
	// DefaultLogMaxSize is the default size of log files.
	DefaultLogMaxSize = 96
	// DefaultLogAge is the default age of the log.
	DefaultLogAge = 30
	// DefaultLogBackups is the default number of log backups.
	DefaultLogBackups = 3
	// DefaultLogCompress is the default value for whether compress log.
	DefaultLogCompress = false
)

var Logger zap.Logger

type LogConfig struct {
	Level zapcore.Level

	LumberjackConfig lumberjack.Logger
}

func NewLogConfig(Level zapcore.Level, LumberjackConfig lumberjack.Logger) LogConfig {
	return LogConfig{
		Level:            Level,
		LumberjackConfig: LumberjackConfig,
	}
}

func InitLog(conf LogConfig) {
	// 编码器配置
	config := zap.NewProductionEncoderConfig()
	// 指定时间编码器
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	// 日志级别用大写
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	// 编码器
	encoder := zapcore.NewConsoleEncoder(config)

	// 配置
	//lj := &lumberjack.Logger{
	//	Filename:   DefaultLogFilePath,
	//	MaxSize:    DefaultLogMaxSize,
	//	MaxBackups: DefaultLogBackups,
	//	MaxAge:     DefaultLogAge,
	//	Compress:   DefaultLogCompress,
	//}

	sync := zapcore.AddSync(&conf.LumberjackConfig)

	// create Logger
	// sync is config for writing to file.
	// zapcore.AddSync(os.Stdout) is config for writing to console.

	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(sync, zapcore.AddSync(os.Stdout)), conf.Level)
	Logger := zap.New(core, zap.AddCaller())

	Logger.Info("Completed init the logger")
}
