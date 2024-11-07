package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"
	"testing"
	"time"
)

type testData struct {
	StuId string `json:"stu_id"`
}

// TestPlaintext
// 阿里云测试通过
func TestPlaintext(t *testing.T) {
	configStruct := &config.Configs{}
	config.Init("./data/config.toml", configStruct)
	logger.Init("TestPlaintext", configStruct.LoggerConfig)

	kafkaConfig := configStruct.Kafka["test-plaintext"]

	//写数据
	writer, err := NewWriter(kafkaConfig, kafkaConfig.Producers[0])
	assert.Nil(t, err)

	data := &testData{StuId: "1"}

	marshalData, err := json.Marshal(data)
	assert.Nil(t, err)

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Value: marshalData,
	})
	assert.Nil(t, err)

	time.Sleep(time.Second)

	//读数据
	reader, err := NewReader(kafkaConfig, kafkaConfig.Consumers[0])
	assert.Nil(t, err)

	message, err := reader.ReadMessage(context.Background())
	assert.Nil(t, err)

	data = &testData{}
	err = json.Unmarshal(message.Value, data)
	assert.Nil(t, err)
}

// TestSaslSslProducer
// 阿里云测试通过
func TestSaslSslProducer(t *testing.T) {
	configStruct := &config.Configs{}
	config.Init("./data/config.toml", configStruct)
	logger.Init("TestSaslSslProducer", configStruct.LoggerConfig)

	kafkaConfig := configStruct.Kafka["test-sasl_ssl"]

	//写数据
	writer, err := NewWriter(kafkaConfig, kafkaConfig.Producers[0])
	assert.Nil(t, err)

	data := &testData{StuId: "1"}

	marshalData, err := json.Marshal(data)
	assert.Nil(t, err)

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Value: marshalData,
	})
	assert.Nil(t, err)

	time.Sleep(time.Second)

	//读数据
	reader, err := NewReader(kafkaConfig, kafkaConfig.Consumers[0])
	assert.Nil(t, err)

	message, err := reader.ReadMessage(context.Background())
	assert.Nil(t, err)

	data = &testData{}
	err = json.Unmarshal(message.Value, data)
	assert.Nil(t, err)
}

func TestSaslPlaintextProducer(t *testing.T) {
	configStruct := &config.Configs{}
	config.Init("./data/config.toml", configStruct)
	logger.Init("TestSaslSslProducer", configStruct.LoggerConfig)

	kafkaConfig := configStruct.Kafka["test-sasl_plaintext"]

	//写数据
	writer, err := NewWriter(kafkaConfig, kafkaConfig.Producers[0])
	assert.Nil(t, err)

	data := &testData{StuId: "1"}

	marshalData, err := json.Marshal(data)
	assert.Nil(t, err)

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Value: marshalData,
	})
	assert.Nil(t, err)

	time.Sleep(time.Second)

	//读数据
	reader, err := NewReader(kafkaConfig, kafkaConfig.Consumers[0])
	assert.Nil(t, err)

	message, err := reader.ReadMessage(context.Background())
	assert.Nil(t, err)

	data = &testData{}
	err = json.Unmarshal(message.Value, data)
	assert.Nil(t, err)
}
