package logger

import (
	"context"
	"fmt"
	"go-web-template/base/common/constant"
	"go-web-template/base/common/utils"
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
	logger *zap.Logger
)

const (
	// defaultPath 日志路径
	defaultPath = "./"
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
func Init(projectName string, loggerConfig config.LoggerConfig) {
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
	//file logger
	logFileCore := zapcore.NewCore(logFileEncoder, output, levelEnable)
	logFileCore = logFileCore.With([]zap.Field{zap.String("projectName", projectName)})
	cores = append(cores, logFileCore)

	//console logger
	if loggerConfig.Console {
		consoleCore := zapcore.NewCore(consoleEncoder, console, levelEnable)
		consoleCore = consoleCore.With([]zap.Field{zap.String("projectName", projectName)})
		cores = append(cores, consoleCore)
	}

	core := zapcore.NewTee(cores...)

	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
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

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func checkNil() {
	if logger == nil {
		fmt.Println("logger is nil, please init logger first")
		panic("logger is nil, please init logger first")
	}
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	checkNil()
	logger.Debug(msg, withCommonField(ctx, fields)...)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	checkNil()
	logger.Info(msg, withCommonField(ctx, fields)...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	checkNil()
	logger.Warn(msg, withCommonField(ctx, fields)...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	checkNil()
	logger.Error(msg, withCommonField(ctx, fields)...)
}

func withCommonField(ctx context.Context, fields []zap.Field) []zap.Field {
	startTime := utils.GetStartTime(ctx)
	span := utils.GetSpan(ctx)

	commonFieldCount := 6
	result := make([]zap.Field, commonFieldCount, len(fields)+commonFieldCount)
	result[0] = zap.String(constant.LoggerKeyType, constant.LoggerTypeLog)
	result[1] = zap.Int64(constant.LoggerKeyTimestampUs, time.Now().UnixMicro())
	result[2] = zap.Int64(constant.LoggerKeyDurationUs, time.Since(startTime).Microseconds())
	result[3] = zap.String(constant.LoggerKeyTraceId, utils.GetTraceId(ctx))
	result[4] = zap.String(constant.LoggerKeyParentSpan, span.GetParentSpan())
	result[5] = zap.String(constant.LoggerKeySpan, span.Span())

	return append(result, fields...)
}

func Flush() error {
	checkNil()
	return logger.Sync()
}
