package signal

import (
	"context"
	"go-web-template/base/lib/logger"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func HandleSignal(ctx context.Context, close func()) {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	sig := <-signalCh

	logger.Info(ctx, "receive shutdown signal", zap.String("signal", sig.String()))

	//处理需要关闭的资源
	close()

	os.Exit(0)
}