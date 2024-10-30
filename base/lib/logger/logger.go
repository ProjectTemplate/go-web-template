package logger

import (
	"fmt"
	"go-web-template/base/lib/config"
	"os"
	"path"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger        *zap.Logger
	loggerSugared *zap.SugaredLogger
)

const (
	// defaultPath 日志路径
	defaultPath = "/var/log/go-web"
	// defaultFileName 日志文件名
	defaultFileName = "server.log"
	// defaultLevel 日志级别
	defaultLevel = zapcore.DebugLevel
	// defaultMaxSize 日志文件最大大小，单位MB
	defaultMaxSize = 100
	// defaultMaxBackups 日志文件最大保存个数
	defaultMaxBackups = 30
	// defaultMaxAge 日志文件最大保存天数
	defaultMaxAge = 28
)

// Init 初始化日志
func Init(loggerConfig config.LoggerConfig) {
	filePath := defaultPath
	fileName := defaultFileName
	level := defaultLevel
	maxSize := defaultMaxSize
	maxBackups := defaultMaxBackups
	maxAge := defaultMaxAge

	if loggerConfig.Path != "" {
		filePath = loggerConfig.Path
	}
	if loggerConfig.FileName != "" {
		fileName = loggerConfig.FileName
	}
	logFile := path.Join(filePath, fileName)
	fmt.Println("log file:", logFile)

	switch strings.ToUpper(loggerConfig.Level) {
	case "DEBUG":
		level = zapcore.DebugLevel
	case "INFO":
		level = zapcore.InfoLevel
	case "WARN":
		level = zapcore.WarnLevel
	case "ERROR":
		level = zapcore.ErrorLevel
	}

	if loggerConfig.MaxSize > 0 {
		maxSize = loggerConfig.MaxSize
	}

	if loggerConfig.MaxBackups > 0 {
		maxBackups = loggerConfig.MaxBackups
	}

	if loggerConfig.MaxAge > 0 {
		maxAge = loggerConfig.MaxAge
	}

	levelEnable := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= level
	})

	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		LocalTime:  true,
		Compress:   false,
	}
	output := zapcore.AddSync(lumberjackLogger)

	console := zapcore.Lock(os.Stdout)
	logFileEncoder := zapcore.NewJSONEncoder(NewEncoderConfig())
	consoleEncoder := zapcore.NewConsoleEncoder(NewEncoderConfig())

	var cores = make([]zapcore.Core, 0, 2)
	cores = append(cores, zapcore.NewCore(logFileEncoder, output, levelEnable))
	if loggerConfig.Console {
		cores = append(cores, zapcore.NewCore(consoleEncoder, console, levelEnable))
	}

	core := zapcore.NewTee(cores...)

	l := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	logger = l
	loggerSugared = l.Sugar()
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
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func Logger() *zap.Logger {
	checkNil()
	return logger
}

func Sugared() *zap.SugaredLogger {
	checkNil()
	return loggerSugared
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func checkNil() {
	if logger == nil {
		fmt.Println("loggerSugared is nil, please init logger first")
		panic("loggerSugared is nil, please init logger first")
	}
}
