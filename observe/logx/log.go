package logx

import (
	"github.com/creasty/defaults"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/dealancer/validate.v2"
	"os"
	"time"
)

func New(cfg Config) (*zap.Logger, error) {
	if err := defaults.Set(&cfg); err != nil {
		return nil, err
	}
	if err := validate.Validate(&cfg); err != nil {
		return nil, err
	}
	encoder := setEncode(cfg)
	write := setWriteSynce(cfg)
	level := setLevel(cfg)

	core := zapcore.NewCore(encoder, write, level)

	ops := []zap.Option{}
	if cfg.Caller {
		ops = append(ops, zap.AddCaller())
	}
	if cfg.Development {
		ops = append(ops, zap.Development())
	}
	return zap.New(core, ops...), nil
}

// 解析器
func setEncode(cfg Config) zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",  //日志时间的key
		LevelKey:       "level", //日志level的key
		NameKey:        "logger",
		CallerKey:      "linenum", //日志产生的文件及其行数的key
		MessageKey:     "msg",     //日志内容的key
		StacktraceKey:  "stacktrace",
		FunctionKey:    "func",                         // 日志函数的key
		LineEnding:     zapcore.DefaultLineEnding,      //回车换行
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器 info debug not Info Debug
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器 包名、文件名、行号
		EncodeName:     zapcore.FullNameEncoder,
	}
	if cfg.HasTimestamp {
		encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			// nanos := t.Unix()
			// sec := float64(nanos) / float64(time.Second)
			enc.AppendInt(int(t.UnixMilli()))
		}
	}
	if cfg.Formatter == JSON {
		return zapcore.NewJSONEncoder(encoderConfig)
	} else {
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
}

// 输出
func setWriteSynce(cfg Config) (write zapcore.WriteSyncer) {
	var hook lumberjack.Logger
	if cfg.OutType == FILE || cfg.OutType == ALL {

		hook = lumberjack.Logger{
			Filename:   cfg.LogName,    // 日志文件路径，默认 os.TempDir()
			MaxSize:    cfg.MaxSize,    // 每个日志文件保存10M，默认 100M
			MaxBackups: cfg.MaxBackNum, // 保留30个备份，默认不限
			MaxAge:     cfg.MaxAge,     // 保留7天，默认不限
			Compress:   cfg.Compress,   // 是否压缩，默认不压缩
		}
	}

	switch cfg.OutType {
	case CONSOLE, "":
		write = zapcore.AddSync(os.Stdout)
	case FILE:
		write = zapcore.AddSync(&hook)
	case ALL:
		write = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook))
	}
	return
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
