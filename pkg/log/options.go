package log

import (
	"flag"
	"github.com/spf13/pflag"
	uzap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"time"
)

type Options struct {
	zap.Options
}

func NewOptions() *Options {
	return &Options{
		Options: zap.Options{
			Development:     false,
			Encoder:         GetEncoder(),
			DestWritter:     GetWriteStdout(),
			Level:           GetLevelEnabler(),
			StacktraceLevel: GetStackLevelEnabler(),
			ZapOpts:         []uzap.Option{uzap.AddCaller()},
		},
	}
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	o.BindFlags(flag.CommandLine)
}

// GetEncoder 自定义的Encoder
func GetEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(
		zapcore.EncoderConfig{
			TimeKey:   "ts",
			LevelKey:  "level",
			NameKey:   "logger",
			CallerKey: "caller_line",
			//FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    cEncodeLevel,
			EncodeTime:     cEncodeTime,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   cEncodeCaller,
		})
}

// GetConsoleEncoder 输出日志到控制台
func GetDevelopmentConsoleEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(uzap.NewDevelopmentEncoderConfig())
}

// GetConsoleEncoder 输出日志到控制台
func GetProductionConsoleEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(uzap.NewProductionEncoderConfig())
}

// GetWriteSyncer 自定义的WriteSyncer
func GetWriteSyncer() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./zap.log",
		MaxSize:    200,
		MaxBackups: 10,
		MaxAge:     30,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// GetWriteSyncer 自定义的WriteSyncer
func GetWriteStdout() zapcore.WriteSyncer {
	return zapcore.AddSync(os.Stdout)
}

// GetLevelEnabler 自定义的LevelEnabler
func GetLevelEnabler() zapcore.Level {
	return zapcore.InfoLevel
}

// GetStackLevelEnabler 自定义的LevelEnabler
func GetStackLevelEnabler() zapcore.Level {
	return zapcore.PanicLevel
}

// cEncodeLevel 自定义日志级别显示
func cEncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(level.CapitalString())
}

// cEncodeTime 自定义时间格式显示
func cEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	logTmFmt := "2006-01-02 15:04:05.000"
	enc.AppendString(t.Format(logTmFmt))
}

// cEncodeCaller 自定义行号显示
func cEncodeCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(caller.TrimmedPath() + "|")
}
