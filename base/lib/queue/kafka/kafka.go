package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"go-web-template/base/lib/logger"
	"go.uber.org/zap"
	"time"

	"github.com/segmentio/kafka-go/sasl/plain"

	"go-web-template/base/common/utils"
	"go-web-template/base/lib/config"
)

var kafkaClusters map[string]*kafkaWriterReader

type kafkaWriterReader struct {
	writers map[string]*kafka.Writer
	readers map[string]*kafka.Reader
}

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

	writer, ok := cluster.writers[writerName]
	if !ok {
		logger.Info(ctx, "GetWriter, writer not found", zap.String("clusterName", clusterName), zap.String("writerName", writerName))
		panic("GetWriter, writer not found, writer name:" + writerName)
	}

	return writer
}

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

	reader, ok := cluster.readers[readerName]
	if !ok {
		logger.Info(ctx, "GetReader, reader not found", zap.String("clusterName", clusterName), zap.String("readerName", readerName))
		panic("GetReader, reader not found, reader name:" + readerName)
	}

	return reader
}

func Init(ctx context.Context, kafkaConfigs map[string]config.Kafka) {
	if len(kafkaConfigs) == 0 {
		logger.Info(ctx, "Init Kafka, configs is empty")
		panic("Init Kafka, configs is empty")
	}

	kafkaClusters = map[string]*kafkaWriterReader{}
	for clusterName, clusterConfig := range kafkaConfigs {
		logger.Info(ctx, "Init Kafka, infos", zap.String("clusterName", clusterName), zap.Any("clusterConfig", clusterConfig))
		writerReaders := &kafkaWriterReader{
			writers: make(map[string]*kafka.Writer),
			readers: make(map[string]*kafka.Reader),
		}

		// 消费者初始化
		for _, consumerConfig := range clusterConfig.Consumers {
			reader, err := newReader(clusterConfig, consumerConfig)
			if err != nil {
				logger.Info(ctx, "Init Kafka, newReader error", zap.String("clusterName", clusterName), zap.Any("consumerConfig", consumerConfig), zap.Error(err))
				panic("Init Kafka, newReader error, cluster name:" + clusterName + "error: " + err.Error())
			}
			writerReaders.readers[consumerConfig.Name] = reader
		}

		// 生产者初始化
		for _, producerConfig := range clusterConfig.Producers {
			writer, err := newWriter(clusterConfig, producerConfig)
			if err != nil {
				logger.Info(ctx, "Init Kafka, newWriter error", zap.String("clusterName", clusterName), zap.Any("producerConfig", producerConfig), zap.Error(err))
				panic("Init Kafka, newWriter error, cluster name:" + clusterName + "error: " + err.Error())
			}
			writerReaders.writers[producerConfig.Name] = writer
		}

		kafkaClusters[clusterName] = writerReaders
	}
}

func newWriter(clusterConfig config.Kafka, producerConfig config.KafkaProducer) (*kafka.Writer, error) {
	ackConfig := kafka.RequireNone
	if producerConfig.AckConfig == "all" {
		ackConfig = kafka.RequireAll
	}
	if producerConfig.AckConfig == "one" {
		ackConfig = kafka.RequireOne
	}

	if clusterConfig.SecurityProtocol == config.SecurityProtocolSaslSsl {
		return newSaslSslProducer(clusterConfig, ackConfig, producerConfig)
	}

	if clusterConfig.SecurityProtocol == config.SecurityProtocolSaslPlaintext {
		return newSaslPlaintextProducer(clusterConfig, ackConfig, producerConfig), nil
	}

	return newPlaintextProducer(clusterConfig, ackConfig, producerConfig), nil
}

func newReader(clusterConfig config.Kafka, consumerConfig config.KafkaConsumer) (*kafka.Reader, error) {
	if clusterConfig.SecurityProtocol == config.SecurityProtocolSaslSsl {
		return newSaslSslReader(clusterConfig, consumerConfig)
	}

	if clusterConfig.SecurityProtocol == config.SecurityProtocolSaslPlaintext {
		return newSaslPlaintextReader(clusterConfig, consumerConfig), nil
	}

	return newPlaintextReader(clusterConfig, consumerConfig), nil
}

func newPlaintextProducer(clusterConfig config.Kafka, ackConfig kafka.RequiredAcks, producerConfig config.KafkaProducer) *kafka.Writer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(clusterConfig.Brokers...),
		Topic:        producerConfig.Topic,
		RequiredAcks: ackConfig,
		Balancer:     &kafka.Hash{},
		Logger:       &kafkaLogger{},
	}
	return w
}

func newSaslSslProducer(clusterConfig config.Kafka, ackConfig kafka.RequiredAcks, producerConfig config.KafkaProducer) (*kafka.Writer, error) {
	tlsConfig, err := utils.NewTlsConfig(clusterConfig.CertData)
	if err != nil {
		return nil, err
	}

	producer := newPlaintextProducer(clusterConfig, ackConfig, producerConfig)
	producer.Transport = &kafka.Transport{
		TLS: tlsConfig,
		SASL: plain.Mechanism{
			Username: clusterConfig.Username,
			Password: clusterConfig.Password,
		},
	}
	return producer, nil
}

func newSaslPlaintextProducer(clusterConfig config.Kafka, ackConfig kafka.RequiredAcks, producerConfig config.KafkaProducer) *kafka.Writer {
	producer := newPlaintextProducer(clusterConfig, ackConfig, producerConfig)
	producer.Transport = &kafka.Transport{
		SASL: plain.Mechanism{
			Username: clusterConfig.Username,
			Password: clusterConfig.Password,
		},
	}
	return producer
}

func newPlaintextReader(clusterConfig config.Kafka, consumerConfig config.KafkaConsumer) *kafka.Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        clusterConfig.Brokers,
		GroupID:        consumerConfig.Group,
		Topic:          consumerConfig.Topic,
		CommitInterval: consumerConfig.CommitInterval,
		Logger:         &kafkaLogger{},
	})
	return reader
}

func newSaslSslReader(clusterConfig config.Kafka, consumerConfig config.KafkaConsumer) (*kafka.Reader, error) {
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

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        clusterConfig.Brokers,
		GroupID:        consumerConfig.Group,
		Topic:          consumerConfig.Topic,
		CommitInterval: consumerConfig.CommitInterval,
		Dialer:         dialer,
	})

	return reader, nil
}

func newSaslPlaintextReader(clusterConfig config.Kafka, consumerConfig config.KafkaConsumer) *kafka.Reader {
	dialer := &kafka.Dialer{
		Timeout:   5 * time.Second,
		DualStack: true,
		SASLMechanism: plain.Mechanism{
			Username: clusterConfig.Username,
			Password: clusterConfig.Password,
		},
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        clusterConfig.Brokers,
		GroupID:        consumerConfig.Group,
		Topic:          consumerConfig.Topic,
		CommitInterval: consumerConfig.CommitInterval,
		Dialer:         dialer,
	})

	return reader
}
