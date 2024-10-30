package logger

import (
	"context"
	"io"
)

var _ io.Writer = (*Proxy)(nil)

// Proxy 日志代理
type Proxy struct {
}

func GetLoggerProxy() *Proxy {
	loggerProxyOnce.Do(func() {
		loggerProxy = &Proxy{}
	})
	return loggerProxy
}

func (p *Proxy) Write(data []byte) (n int, err error) {
	Info(context.Background(), string(data))
	return len(data), nil
}
