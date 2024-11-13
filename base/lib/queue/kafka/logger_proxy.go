package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go-web-template/base/lib/logger"
	"go.uber.org/zap"
)

var _ kafka.Logger = (*kafkaLogger)(nil)

type kafkaLogger struct {
}

func (k *kafkaLogger) Printf(s string, i ...interface{}) {
	message := fmt.Sprintf(s, i...)
	logger.Info(context.Background(), message, zap.String("tag", "kafka"))
}
