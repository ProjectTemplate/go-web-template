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
	logger        *zap.Logger
	loggerSugared *zap.SugaredLogger
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
	newCore := zapcore.NewCore(logFileEncoder, output, levelEnable)
	newCore = newCore.With([]zap.Field{zap.String("projectName", projectName)})
	cores = append(cores, newCore)
	if loggerConfig.Console {
		consoleCore := zapcore.NewCore(consoleEncoder, console, levelEnable)
		consoleCore = consoleCore.With([]zap.Field{zap.String("projectName", projectName)})
		cores = append(cores, consoleCore)
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

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func checkNil() {
	if logger == nil {
		fmt.Println("loggerSugared is nil, please init logger first")
		panic("loggerSugared is nil, please init logger first")
	}
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	checkNil()
	mergerFields := append(commonLoggerFields(ctx), fields...)
	logger.Debug(msg, mergerFields...)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	checkNil()
	mergerFields := append(commonLoggerFields(ctx), fields...)
	logger.Info(msg, mergerFields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	checkNil()
	mergerFields := append(commonLoggerFields(ctx), fields...)
	logger.Warn(msg, mergerFields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	checkNil()
	mergerFields := append(commonLoggerFields(ctx), fields...)
	logger.Error(msg, mergerFields...)
}

func SDebugF(ctx context.Context, msg string, args ...interface{}) {
	checkNil()
	message := formatMessage(msg, args...)
	loggerSugared.Debugw(message, commonLoggerKeyValues(ctx)...)
}

func SInfoF(ctx context.Context, msg string, args ...interface{}) {
	checkNil()
	message := formatMessage(msg, args...)
	loggerSugared.Infow(message, commonLoggerKeyValues(ctx)...)
}

func SWarnF(ctx context.Context, msg string, args ...interface{}) {
	checkNil()
	message := formatMessage(msg, args...)
	loggerSugared.Warnw(message, commonLoggerKeyValues(ctx)...)
}

func SErrorF(ctx context.Context, msg string, args ...interface{}) {
	checkNil()
	message := formatMessage(msg, args...)
	loggerSugared.Errorw(message, commonLoggerKeyValues(ctx)...)
}

func formatMessage(template string, args ...interface{}) string {
	if len(args) > 0 {
		return fmt.Sprintf(template, args...)
	}
	return template
}

func commonLoggerKeyValues(ctx context.Context) []interface{} {
	spanStr := ""
	span := utils.GetSpan(ctx)
	if span != nil {
		spanStr = span.IncreaseAndGet()
	}

	return []interface{}{
		constant.ContextKeyDomain, utils.GetDomain(ctx),
		constant.HeaderKeyTraceId, utils.GetTraceId(ctx),
		constant.ContextKeySpan, spanStr,
	}
}

func commonLoggerFields(ctx context.Context) []zap.Field {
	loggerKeyValues := commonLoggerKeyValues(ctx)

	result := make([]zap.Field, 0, len(loggerKeyValues)/2)
	for i := 0; i < len(loggerKeyValues); i += 2 {
		result = append(result, zap.String(loggerKeyValues[i].(string), loggerKeyValues[i+1].(string)))
	}

	return result
}
