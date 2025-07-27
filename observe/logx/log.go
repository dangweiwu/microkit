package logx

import (
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 添加包级别的原子级别变量
var atomicLevel zap.AtomicLevel

func New(cfg Config) (*zap.Logger, error) {
	atomicLevel = zap.NewAtomicLevelAt(setLevel(cfg))
	zapCfg := zap.Config{
		Level:         atomicLevel,
		Development:   cfg.Development,
		DisableCaller: !cfg.Caller,
		Encoding:      cfg.Formatter,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "linenum",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			FunctionKey:    "func",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
	}

	if cfg.HasTimestamp {
		zapCfg.EncoderConfig.EncodeTime = customTimeEncoder
	}

	// 使用getFileWriter创建带轮转功能的写入器
	var writer zapcore.WriteSyncer
	fileWriter := getFileWriter(cfg)
	switch cfg.OutType {
	case CONSOLE:
		writer = zapcore.AddSync(os.Stdout)
	case FILE:
		writer = zapcore.AddSync(fileWriter)
	default:
		writer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter))
	}

	// 创建编码器
	var encoder zapcore.Encoder
	switch zapCfg.Encoding {
	case JSON:
		encoder = zapcore.NewJSONEncoder(zapCfg.EncoderConfig)
	default:
		encoder = zapcore.NewConsoleEncoder(zapCfg.EncoderConfig)
	}

	// 创建核心并应用轮转写入器
	core := zapcore.NewCore(encoder, writer, zapCfg.Level)

	// 使用自定义核心构建logger

	ops := []zap.Option{zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return core
	})}

	return zapCfg.Build(ops...)

}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt(int(t.UnixMilli()))
}

// getFileWriter 创建文件写入器
func getFileWriter(cfg Config) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   cfg.LogName,    // 日志文件路径，默认 os.TempDir()
		MaxSize:    cfg.MaxSize,    // 每个日志文件保存10M，默认 100M
		MaxBackups: cfg.MaxBackNum, // 保留30个备份，默认不限
		MaxAge:     cfg.MaxAge,     // 保留7天，默认不限
		Compress:   cfg.Compress,   // 是否压缩，默认不压缩
	}

	return zapcore.AddSync(lumberJackLogger)
}

// 类型
func setLevel(cfg Config) (level zapcore.Level) {
	switch cfg.Level {
	case INFO:
		level = zapcore.InfoLevel
	case DEBUG:
		level = zapcore.DebugLevel
	case WARN:
		level = zapcore.WarnLevel
	case ERROR:
		level = zapcore.ErrorLevel
	case PANIC: //会跳出来
		level = zapcore.DPanicLevel
	default:
		level = zapcore.DebugLevel
	}
	return
}

// 添加动态设置级别方法
func SetLevel(level string) {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}
	atomicLevel.SetLevel(zapLevel)
}
