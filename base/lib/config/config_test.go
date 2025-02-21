package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
	var conf = &Configs{}
	Init("./data/config.toml", conf)

	assert.NotNil(t, conf)
	assert.NotNil(t, conf.App)
	assert.NotNil(t, conf.Server)
	assert.NotNil(t, conf.LoggerConfig)

	assert.Equal(t, "test", conf.App.Name)

	assert.Equal(t, "127.0.0.1:8080", conf.Server.Address)
	assert.Equal(t, true, conf.Server.Debug)

	assert.Equal(t, true, conf.LoggerConfig.Console)
	assert.Equal(t, "DEBUG", conf.LoggerConfig.Level)
	assert.Equal(t, "./", conf.LoggerConfig.Path)
	assert.Equal(t, "server.log", conf.LoggerConfig.FileName)
	assert.Equal(t, 100, conf.LoggerConfig.MaxSize)
	assert.Equal(t, 30, conf.LoggerConfig.MaxBackups)
	assert.Equal(t, 15, conf.LoggerConfig.MaxAge)

	assert.Equal(t, time.Millisecond*500, conf.FastHttp.ReadTimeOut)
	assert.Equal(t, time.Millisecond*500, conf.FastHttp.WriteTimeOut)
	assert.Equal(t, time.Hour, conf.FastHttp.MaxIdleConnDuration)
	assert.Equal(t, 512, conf.FastHttp.MaxConnsPerHost)
	assert.Equal(t, 2, conf.FastHttp.RetryTimes)

	assert.NotEmpty(t, conf.Mysql)
	assert.Equal(t, 1, len(conf.Mysql))
	assert.Equal(t, "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&loc=Local&parseTime=True", conf.Mysql["test"].DSN[0])
	assert.Equal(t, 50, conf.Mysql["test"].MaxOpenConnections)
	assert.Equal(t, 25, conf.Mysql["test"].MaxIdleConnections)
	assert.Equal(t, time.Hour, conf.Mysql["test"].MaxLifeTime)
	assert.Equal(t, time.Minute*10, conf.Mysql["test"].MaxIdleTime)
	assert.Equal(t, true, conf.Mysql["test"].ShowLog)
	assert.Equal(t, time.Second, conf.Mysql["test"].SlowThreshold)

	assert.NotEmpty(t, conf.Nacos)
	assert.Equal(t, 2, len(conf.Nacos))
	assert.Equal(t, "all", conf.Nacos["test"].ClientType)
	assert.Equal(t, time.Second*10, conf.Nacos["test"].TimeOut)
	assert.Equal(t, "info", conf.Nacos["test"].LogLevel)
	assert.Equal(t, true, conf.Nacos["test"].AppendToStdout)
	assert.Equal(t, "test", conf.Nacos["test"].Namespace)
	assert.NotEmpty(t, conf.Nacos["test"].Servers)
	assert.Equal(t, 2, len(conf.Nacos["test"].Servers))

	assert.Equal(t, "all", conf.Nacos["test1"].ClientType)
	assert.Equal(t, time.Second*10, conf.Nacos["test1"].TimeOut)
	assert.Equal(t, "info", conf.Nacos["test1"].LogLevel)
	assert.Equal(t, true, conf.Nacos["test1"].AppendToStdout)
	assert.Equal(t, "test", conf.Nacos["test1"].Namespace)
	assert.NotEmpty(t, conf.Nacos["test1"].Servers)
	assert.Equal(t, 2, len(conf.Nacos["test1"].Servers))

	assert.NotEmpty(t, conf.Redis)
	assert.Equal(t, 2, len(conf.Redis))
	assert.Equal(t, "127.0.0.1:6379", conf.Redis["test"].Addr)
	assert.Equal(t, "", conf.Redis["test"].Password)
	assert.Equal(t, 0, conf.Redis["test"].DB)
	assert.Equal(t, "127.0.0.1:6379", conf.Redis["test1"].Addr)
	assert.Equal(t, "", conf.Redis["test1"].Password)
	assert.Equal(t, 1, conf.Redis["test1"].DB)

	assert.NotEmpty(t, conf.Kafka)
	assert.Equal(t, 3, len(conf.Kafka))
	assert.NotEmpty(t, conf.Kafka["test-plaintext"].Brokers)
	assert.Equal(t, "plaintext", conf.Kafka["test-plaintext"].SecurityProtocol)
	assert.Equal(t, 1, len(conf.Kafka["test-plaintext"].Writers))
	assert.Equal(t, 1, len(conf.Kafka["test-plaintext"].Readers))
	assert.Equal(t, "test", conf.Kafka["test-plaintext"].Writers[0].Name)
	assert.Equal(t, "test", conf.Kafka["test-plaintext"].Readers[0].Name)

	assert.NotEmpty(t, conf.Kafka["test-sasl_ssl"].Brokers)
	assert.Equal(t, "sasl_ssl", conf.Kafka["test-sasl_ssl"].SecurityProtocol)
	assert.Equal(t, 1, len(conf.Kafka["test-sasl_ssl"].Writers))
	assert.Equal(t, 1, len(conf.Kafka["test-sasl_ssl"].Readers))
	assert.Equal(t, "test", conf.Kafka["test-sasl_ssl"].Writers[0].Name)
	assert.Equal(t, "test", conf.Kafka["test-sasl_ssl"].Readers[0].Name)

	assert.NotEmpty(t, conf.Kafka["test-sasl_plaintext"].Brokers)
	assert.Equal(t, "sasl_plaintext", conf.Kafka["test-sasl_plaintext"].SecurityProtocol)
	assert.Equal(t, 1, len(conf.Kafka["test-sasl_plaintext"].Writers))
	assert.Equal(t, 1, len(conf.Kafka["test-sasl_plaintext"].Readers))
	assert.Equal(t, "test", conf.Kafka["test-sasl_plaintext"].Writers[0].Name)
	assert.Equal(t, "test", conf.Kafka["test-sasl_plaintext"].Readers[0].Name)

	assert.Equal(t, "go-web-template", conf.Otel.Trace.ScopeName)
	assert.Equal(t, "127.0.0.1:4318", conf.Otel.Trace.Endpoint)
	assert.Equal(t, true, conf.Otel.Trace.Insecure)
	assert.Equal(t, "service_namespace_test", conf.Otel.Trace.ServiceNamespace)
	assert.Equal(t, "service_test", conf.Otel.Trace.ServiceName)
	assert.Equal(t, "instance_test", conf.Otel.Trace.ServiceInstanceID)
	assert.Equal(t, "v1.0.0", conf.Otel.Trace.ServiceVersion)

	assert.Equal(t, "cluster_test", conf.Otel.Metric.Cluster)
	assert.Equal(t, "service_test", conf.Otel.Metric.ServiceName)
	assert.Equal(t, true, conf.Otel.Metric.Insecure)
	assert.Equal(t, "127.0.0.1:4318", conf.Otel.Metric.Endpoint)
}
