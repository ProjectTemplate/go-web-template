package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go-web-template/base/lib/logger"
	"go.uber.org/zap"
)

var _ kafka.Logger = (*kafkaErrorLogger)(nil)

type kafkaErrorLogger struct {
}

func (k *kafkaErrorLogger) Printf(s string, i ...interface{}) {
	message := fmt.Sprintf(s, i...)
	logger.Error(context.Background(), message, zap.String("tag", "kafka"))
}
