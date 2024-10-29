package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.SugaredLogger

// Init 初始化日志
func Init() {
	GTEError := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	GTEDebug := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl > zapcore.DebugLevel
	})

	infoOutput := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "log-info.log",
		MaxSize:    10,
		MaxBackups: 30,
		MaxAge:     28,
	})

	errorOutput := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "log-error.log",
		MaxSize:    10,
		MaxBackups: 30,
		MaxAge:     28,
	})

	console := zapcore.Lock(os.Stdout)
	prodEncoder := zapcore.NewJSONEncoder(NewEncoderConfig())
	devEncoder := zapcore.NewConsoleEncoder(NewEncoderConfig())

	core := zapcore.NewTee(
		zapcore.NewCore(devEncoder, console, GTEDebug),
		zapcore.NewCore(prodEncoder, infoOutput, GTEDebug),
		zapcore.NewCore(prodEncoder, errorOutput, GTEError),
	)

	l := zap.New(core, zap.AddCaller())
	logger = l.Sugar()
}

func NewEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func New() *zap.SugaredLogger {
	return logger
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}
