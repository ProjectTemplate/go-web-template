package metric

import (
	"context"
	"testing"
	"time"

	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"

	"github.com/stretchr/testify/assert"
)

func TestOtelMetric(t *testing.T) {

	configStruct := &config.Configs{}
	config.Init("./data/config.toml", configStruct)
	logger.Init("TestOtelMetric", configStruct.LoggerConfig)

	ctx := context.Background()

	Init(ctx, configStruct.Otel.Metric)

	err := RecordGoRoutineCount(configStruct.Otel.Metric)
	assert.Nil(t, err)

	err = RecordMemoryInfo(configStruct.Otel.Metric)
	assert.Nil(t, err)

	err = RecordLoadInfo(configStruct.Otel.Metric)
	assert.Nil(t, err)

	err = RecordCpuInfo(configStruct.Otel.Metric)
	assert.Nil(t, err)

	go func() {
		for {
			time.Sleep(time.Millisecond * 10)
			go func() {
				time.Sleep(time.Minute)
			}()
		}
	}()

	time.Sleep(time.Second * 10000)
}
