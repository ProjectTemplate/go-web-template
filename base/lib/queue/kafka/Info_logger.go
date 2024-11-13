package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go-web-template/base/lib/logger"
	"go.uber.org/zap"
)

var _ kafka.Logger = (*kafkaInfoLogger)(nil)

type kafkaInfoLogger struct {
}

func (k *kafkaInfoLogger) Printf(s string, i ...interface{}) {
	message := fmt.Sprintf(s, i...)
	logger.Info(context.Background(), message, zap.String("tag", "kafka"))
}
