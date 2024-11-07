package logger

import (
	"context"
	"io"
)

var _ io.Writer = (*Writer)(nil)

// Writer 日志代理
type Writer struct {
}

func GetWriter() io.Writer {
	loggerProxyOnce.Do(func() {
		loggerProxy = &Writer{}
	})
	return loggerProxy
}

func (p *Writer) Write(data []byte) (n int, err error) {
	Info(context.Background(), string(data))
	return len(data), nil
}

type KafkaLogger struct {
}
