package kafka

import (
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
	"go-web-template/base/lib/logger"
	"go.uber.org/zap"
	"sync"
	"time"

	"github.com/segmentio/kafka-go/sasl/plain"

	"go-web-template/base/common/utils"
	"go-web-template/base/lib/config"
)

var kafkaClusters map[string]*kafkaWriterReader

type kafkaWriterReader struct {
	writers map[string]*kafka.Writer
	// 创建 kafka.Reader 的时候会直接启动 Reader 和服务器建立连接，这里存储配置，调用GetReader方法的时候再创建
	readerConfigs map[string]*kafka.ReaderConfig
	// 用于存储已经创建的 Reader，避免重复创建
	readers sync.Map
}

// GetWriter 根据集群名字和writer名字获取 [kafka.Writer]，如果不存在会panic
func GetWriter(ctx context.Context, clusterName string, writerName string) *kafka.Writer {
	if kafkaClusters == nil {
		logger.Info(ctx, "GetWriter, kafkaClusters is nil")
		panic("GetWriter, kafkaClusters is nil")
	}

	cluster, ok := kafkaClusters[clusterName]
	if !ok {
		logger.Info(ctx, "GetWriter, cluster not found", zap.String("clusterName", clusterName))
		panic("GetWriter, cluster not found, cluster name:" + clusterName)
	}

	//获取writer
	writer, ok := cluster.writers[writerName]
	if !ok {
		logger.Info(ctx, "GetWriter, writer not found", zap.String("clusterName", clusterName), zap.String("writerName", writerName))
		panic("GetWriter, writer not found, writer name:" + writerName)
	}

	return writer
}

// GetReader 根据集群名字和reader名字获取 [kafka.Reader]，如果不存在会panic
func GetReader(ctx context.Context, clusterName string, readerName string) *kafka.Reader {
	if kafkaClusters == nil {
		logger.Info(ctx, "GetReader, kafkaClusters is nil")
		panic("GetReader, kafkaClusters is nil")
	}

	cluster, ok := kafkaClusters[clusterName]
	if !ok {
		logger.Info(ctx, "GetReader, cluster not found", zap.String("clusterName", clusterName))
		panic("GetReader, cluster not found, cluster name:" + clusterName)
	}

	// 从缓存中获取，如果已经初始化过则直接返回
	value, ok := cluster.readers.Load(readerName)
	if ok {
		return value.(*kafka.Reader)
	}

	//获取配置
	readerConfig, ok := cluster.readerConfigs[readerName]
	if !ok {
		logger.Info(ctx, "GetReader, reader not found", zap.String("clusterName", clusterName), zap.String("readerName", readerName))
		panic("GetReader, reader not found, reader name:" + readerName)
	}

	//初始化
	reader := kafka.NewReader(*readerConfig)

	//缓存
	cluster.readers.Store(readerName, reader)

	return reader
}

// Init 初始化kafka
// 支持三种连接方式，分别是 plaintext、sasl_ssl和sasl_plaintext
//
// plaintext 不需要用户密码
//
// sasl_plaintext 需要用户名和密码
//
// sasl_ssl 需要证书、用户名和密码
//
// 详细配置信息可参考配置文件 [./data/config.toml]
func Init(ctx context.Context, kafkaConfigs map[string]config.Kafka) {
	if len(kafkaConfigs) == 0 {
		logger.Info(ctx, "Init Kafka, configs is empty")
		panic("Init Kafka, configs is empty")
	}

	kafkaClusters = map[string]*kafkaWriterReader{}
	for clusterName, clusterConfig := range kafkaConfigs {
		logger.Info(ctx, "Init Kafka, infos", zap.String("clusterName", clusterName), zap.Any("clusterConfig", clusterConfig))
		writerReaders := &kafkaWriterReader{
			writers:       make(map[string]*kafka.Writer),
			readerConfigs: make(map[string]*kafka.ReaderConfig),
			readers:       sync.Map{},
		}

		// 消费者初始化
		for _, readerConfig := range clusterConfig.Readers {
			reader, err := newReader(clusterConfig, readerConfig)
			if err != nil {
				logger.Info(ctx, "Init Kafka, newReader error", zap.String("clusterName", clusterName), zap.Any("readerConfig", readerConfig), zap.Error(err))
				panic("Init Kafka, newReader error, cluster name:" + clusterName + "error: " + err.Error())
			}
			writerReaders.readerConfigs[readerConfig.Name] = reader
		}

		// 生产者初始化
		for _, writerConfig := range clusterConfig.Writers {
			writer, err := newWriter(clusterConfig, writerConfig)
			if err != nil {
				logger.Info(ctx, "Init Kafka, newWriter error", zap.String("clusterName", clusterName), zap.Any("writerConfig", writerConfig), zap.Error(err))
				panic("Init Kafka, newWriter error, cluster name:" + clusterName + "error: " + err.Error())
			}
			writerReaders.writers[writerConfig.Name] = writer
		}

		kafkaClusters[clusterName] = writerReaders
	}
}

func newWriter(clusterConfig config.Kafka, writerConfig config.KafkaWriter) (*kafka.Writer, error) {
	ackConfig := kafka.RequireNone
	if writerConfig.AckConfig == "all" {
		ackConfig = kafka.RequireAll
	}
	if writerConfig.AckConfig == "one" {
		ackConfig = kafka.RequireOne
	}

	if clusterConfig.SecurityProtocol == config.SecurityProtocolSaslPlaintext {
		return newSaslPlaintextWriter(clusterConfig, ackConfig, writerConfig), nil
	}

	if clusterConfig.SecurityProtocol == config.SecurityProtocolSaslSsl {
		return newSaslSslWriter(clusterConfig, ackConfig, writerConfig)
	}

	if clusterConfig.SecurityProtocol == config.SecurityProtocolPlaintext {
		return newPlaintextWriter(clusterConfig, ackConfig, writerConfig), nil
	}

	return nil, errors.New("unsupported security protocol, protocol:" + clusterConfig.SecurityProtocol)
}

func newReader(clusterConfig config.Kafka, readerConfig config.KafkaReader) (*kafka.ReaderConfig, error) {
	if clusterConfig.SecurityProtocol == config.SecurityProtocolSaslSsl {
		return newSaslSslReader(clusterConfig, readerConfig)
	}

	if clusterConfig.SecurityProtocol == config.SecurityProtocolSaslPlaintext {
		return newSaslPlaintextReader(clusterConfig, readerConfig), nil
	}

	return newPlaintextReader(clusterConfig, readerConfig), nil
}

func newPlaintextWriter(clusterConfig config.Kafka, ackConfig kafka.RequiredAcks, writerConfig config.KafkaWriter) *kafka.Writer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(clusterConfig.Brokers...),
		Topic:        writerConfig.Topic,
		RequiredAcks: ackConfig,
		Balancer:     &kafka.Hash{},
		Logger:       &kafkaInfoLogger{},
		ErrorLogger:  &kafkaErrorLogger{},
	}
	return w
}

func newSaslPlaintextWriter(clusterConfig config.Kafka, ackConfig kafka.RequiredAcks, writerConfig config.KafkaWriter) *kafka.Writer {
	writer := newPlaintextWriter(clusterConfig, ackConfig, writerConfig)
	writer.Transport = &kafka.Transport{
		SASL: plain.Mechanism{
			Username: clusterConfig.Username,
			Password: clusterConfig.Password,
		},
	}
	return writer
}

func newSaslSslWriter(clusterConfig config.Kafka, ackConfig kafka.RequiredAcks, writerConfig config.KafkaWriter) (*kafka.Writer, error) {
	tlsConfig, err := utils.NewTlsConfig(clusterConfig.CertData)
	if err != nil {
		return nil, err
	}

	writer := newPlaintextWriter(clusterConfig, ackConfig, writerConfig)
	writer.Transport = &kafka.Transport{
		TLS: tlsConfig,
		SASL: plain.Mechanism{
			Username: clusterConfig.Username,
			Password: clusterConfig.Password,
		},
	}
	return writer, nil
}

func newReaderConfig(clusterConfig config.Kafka, readerConfig config.KafkaReader) *kafka.ReaderConfig {
	return &kafka.ReaderConfig{
		Brokers:        clusterConfig.Brokers,
		GroupID:        readerConfig.Group,
		Topic:          readerConfig.Topic,
		CommitInterval: readerConfig.CommitInterval,
		Logger:         &kafkaInfoLogger{},
		ErrorLogger:    &kafkaErrorLogger{},
	}
}

func newPlaintextReader(clusterConfig config.Kafka, readerConfig config.KafkaReader) *kafka.ReaderConfig {
	return newReaderConfig(clusterConfig, readerConfig)
}

func newSaslPlaintextReader(clusterConfig config.Kafka, readerConfig config.KafkaReader) *kafka.ReaderConfig {
	dialer := &kafka.Dialer{
		Timeout:   5 * time.Second,
		DualStack: true,
		SASLMechanism: plain.Mechanism{
			Username: clusterConfig.Username,
			Password: clusterConfig.Password,
		},
	}

	result := newReaderConfig(clusterConfig, readerConfig)
	result.Dialer = dialer

	return result
}

func newSaslSslReader(clusterConfig config.Kafka, readerConfig config.KafkaReader) (*kafka.ReaderConfig, error) {
	tlsConfig, err := utils.NewTlsConfig(clusterConfig.CertData)
	if err != nil {
		return nil, err
	}

	dialer := &kafka.Dialer{
		Timeout:   5 * time.Second,
		DualStack: true,
		TLS:       tlsConfig,
		SASLMechanism: plain.Mechanism{
			Username: clusterConfig.Username,
			Password: clusterConfig.Password,
		},
	}

	result := newReaderConfig(clusterConfig, readerConfig)
	result.Dialer = dialer

	return result, nil
}
