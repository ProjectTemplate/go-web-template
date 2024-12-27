package kafka

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"

	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"
)

type testData struct {
	StuId string `json:"stu_id"`
}

//go:generate go test -v -run="^TestPlaintext$" .
//go:generate go test -v -run="^TestSaslSslWriter$" .
//go:generate go test -v -run="^TestSaslPlaintextWriter$" .

// TestPlaintext
// 阿里云测试通过
func TestPlaintext(t *testing.T) {
	configStruct := &config.Configs{}
	config.Init("./data/config.toml", configStruct)
	logger.Init("TestPlaintext", configStruct.LoggerConfig)

	background := context.Background()
	Init(background, configStruct.Kafka)

	//写数据
	writer := GetWriter(background, "test-plaintext", "test")

	data := &testData{StuId: "1"}
	marshalData, err := json.Marshal(data)
	assert.Nil(t, err)
	err = writer.WriteMessages(context.Background(), kafka.Message{
		Value: marshalData,
	})
	assert.Nil(t, err)

	time.Sleep(time.Second)

	//读数据
	reader := GetReader(background, "test-plaintext", "test")
	assert.Nil(t, err)

	message, err := reader.ReadMessage(context.Background())
	assert.Nil(t, err)

	data = &testData{}
	err = json.Unmarshal(message.Value, data)
	assert.Nil(t, err)
}

// TestSaslPlaintextWriter
// 腾讯云测试通过
func TestSaslPlaintextWriter(t *testing.T) {
	configStruct := &config.Configs{}
	config.Init("./data/config.toml", configStruct)
	logger.Init("TestSaslPlaintextWriter", configStruct.LoggerConfig)
	background := context.Background()

	Init(background, configStruct.Kafka)

	//写数据
	writer := GetWriter(context.Background(), "test-sasl_plaintext", "test")

	data := &testData{StuId: "1"}
	marshalData, err := json.Marshal(data)
	assert.Nil(t, err)

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Value: marshalData,
	})
	assert.Nil(t, err)

	time.Sleep(time.Second)

	//读数据
	reader := GetReader(context.Background(), "test-sasl_plaintext", "test")
	message, err := reader.ReadMessage(context.Background())
	assert.Nil(t, err)

	data = &testData{}
	err = json.Unmarshal(message.Value, data)
	assert.Nil(t, err)
}

// TestSaslSslWriter
// 阿里云测试通过
// 腾讯云测试通过
func TestSaslSslWriter(t *testing.T) {
	configStruct := &config.Configs{}
	config.Init("./data/config.toml", configStruct)
	logger.Init("TestSaslSslWriter", configStruct.LoggerConfig)
	background := context.Background()

	Init(background, configStruct.Kafka)

	//写数据
	writer := GetWriter(context.Background(), "test-sasl_ssl", "test")

	data := &testData{StuId: "1"}
	marshalData, err := json.Marshal(data)
	assert.Nil(t, err)
	err = writer.WriteMessages(context.Background(), kafka.Message{
		Value: marshalData,
	})
	assert.Nil(t, err)

	time.Sleep(time.Second)

	//读数据
	reader := GetReader(context.Background(), "test-sasl_ssl", "test")

	message, err := reader.ReadMessage(context.Background())
	assert.Nil(t, err)

	data = &testData{}
	err = json.Unmarshal(message.Value, data)
	assert.Nil(t, err)
}

func TestInit(t *testing.T) {
	configStruct := &config.Configs{}
	config.Init("./data/config.toml", configStruct)
	logger.Init("TestInit", configStruct.LoggerConfig)

	background := context.Background()
	Init(background, configStruct.Kafka)

	writer := GetWriter(background, "test-sasl_plaintext", "test")
	assert.NotNil(t, writer)

	reader := GetReader(background, "test-sasl_plaintext", "test")
	assert.NotNil(t, reader)

	reader = GetReader(background, "test-sasl_plaintext", "test")
	assert.NotNil(t, reader)
	time.Sleep(time.Second * 5)
}
