package kafka

import (
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"go-web-template/base/common/utils"
	"go-web-template/base/lib/config"
	"time"
)

func NewWriter(clusterConfig config.Kafka, producerConfig config.KafkaProducer) (*kafka.Writer, error) {
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

func NewReader(clusterConfig config.Kafka, consumerConfig config.KafkaConsumer) (*kafka.Reader, error) {
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
