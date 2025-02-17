package config

import (
	"time"
)

const (
	// SecurityProtocolPlaintext vpc 不需用户名密码
	SecurityProtocolPlaintext = "plaintext"
	// SecurityProtocolSaslSsl ssl 需要证书和用户名密码
	SecurityProtocolSaslSsl = "sasl_ssl"
	// SecurityProtocolSaslPlaintext 只需要用户名、密码
	SecurityProtocolSaslPlaintext = "sasl_plaintext"
)

// Configs 配置信息
type Configs struct {
	App          App              `mapstructure:"app"`
	Server       Server           `mapstructure:"server"`
	LoggerConfig LoggerConfig     `mapstructure:"log"`
	FastHttp     FastHttp         `mapstructure:"fast_http"`
	Mysql        map[string]MySQL `mapstructure:"mysql"`
	Nacos        map[string]Nacos `mapstructure:"nacos"`
	Redis        map[string]Redis `mapstructure:"redis"`
	Kafka        map[string]Kafka `mapstructure:"kafka"`
	Otel         Otel             `mapstructure:"otel"`
}

type App struct {
	Name string `mapstructure:"name"`
}

// Server 服务器配置
type Server struct {
	Address string `mapstructure:"address"`
	Debug   bool   `mapstructure:"debug"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	// Level 日志级别 zapcore.Level
	Level string `mapstructure:"level"`
	// Console 是否输出到控制台
	Console bool `mapstructure:"console"`
	// Path 日志文件路径
	Path string `mapstructure:"path"`
	// FileName 日志文件名
	FileName string `mapstructure:"file_name"`

	// MaxSize 日志文件最大大小，单位MB
	MaxSize int `mapstructure:"max_size"`
	// MaxAge 日志文件最大个数
	MaxBackups int `mapstructure:"max_backups"`
	// MaxAge 日志文件最大保存天数
	MaxAge int `mapstructure:"max_age"`
}

// FastHttp 配置
type FastHttp struct {
	// ReadTimeOut 从读缓冲区读取数据的超时时间，如果在调用的时候指定超时时间，则最短的一个会生效
	ReadTimeOut time.Duration `mapstructure:"read_time_out"`
	// WriteTimeOut 写入响应数据的超时时间，如果在调用的时候指定超时时间，则最短的一个会生效
	WriteTimeOut        time.Duration `mapstructure:"write_time_out"`
	MaxIdleConnDuration time.Duration `mapstructure:"max_idle_conn_duration"`
	MaxConnsPerHost     int           `mapstructure:"max_conns_per_host"`
	RetryTimes          int           `mapstructure:"retry_times"`
}

// MySQL 数据库配置
type MySQL struct {
	DSN                []string      `mapstructure:"dsn"`
	MaxOpenConnections int           `mapstructure:"max_open_connections"`
	MaxIdleConnections int           `mapstructure:"max_idle_connections"`
	MaxLifeTime        time.Duration `mapstructure:"max_life_time"`
	MaxIdleTime        time.Duration `mapstructure:"max_idle_time"`
	ShowLog            bool          `mapstructure:"show_log"`
	SlowThreshold      time.Duration `mapstructure:"slow_threshold"`
}

// Redis Redis配置
type Redis struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// Nacos Nacos配置
type Nacos struct {
	// Type 客户端类型 config 配置客户端，naming 注册中心客户端，all 配置客户端和注册中心客户端
	ClientType     string            `mapstructure:"client_type"`
	TimeOut        time.Duration     `mapstructure:"time_out"`
	LogLevel       string            `mapstructure:"log_level"`
	AppendToStdout bool              `mapstructure:"append_to_stdout"`
	Namespace      string            `mapstructure:"namespace"`
	Servers        []NacosServerConf `mapstructure:"servers"`
}

// NacosServerConf NacosConf 连接配置
type NacosServerConf struct {
	IpAddr string `mapstructure:"ip_addr"`
	Port   uint64 `mapstructure:"port"`
}

// Kafka kafka配置
type Kafka struct {
	Brokers          []string      `mapstructure:"brokers"`
	SecurityProtocol string        `mapstructure:"security_protocol"`
	Username         string        `mapstructure:"username"`
	Password         string        `mapstructure:"password"`
	CertData         string        `mapstructure:"cert_data"`
	Readers          []KafkaReader `mapstructure:"readers"`
	Writers          []KafkaWriter `mapstructure:"writers"`
}

// KafkaReader kafka消费者配置
type KafkaReader struct {
	Name  string `mapstructure:"name"`
	Topic string `mapstructure:"topic"`
	Group string `mapstructure:"group"`
	// CommitInterval 0 表示同步提交
	CommitInterval time.Duration `mapstructure:"commit_interval"`
}

// KafkaWriter kafka生产者配置
type KafkaWriter struct {
	Name  string `mapstructure:"name"`
	Topic string `mapstructure:"topic"`

	// AckConfig
	//
	// none 不等待确认
	//
	// one 等待leader确认
	//
	// all 等待所有在ISR里面的节点确认
	AckConfig string `mapstructure:"ack_config"`
}

type Otel struct {
	Trace Trace `mapstructure:"trace"`
}

type Trace struct {
	// Endpoint 服务端地址 such as "127.0.0.1:4318"
	Endpoint string `mapstructure:"endpoint"`

	// Insecure 是否忽略证书
	Insecure bool `mapstructure:"insecure"`

	// ScopeName 命名空间名字
	ScopeName string `mapstructure:"scope_name"`
	// ServiceNamespace 服务命名空间
	ServiceNamespace string `mapstructure:"service_namespace"`
	// ServiceName 服务名
	ServiceName string `mapstructure:"service_name"`
	// ServiceInstanceID 服务示例ID
	ServiceInstanceID string `mapstructure:"service_instance_id"`
	// ServiceVersion 服务版本
	ServiceVersion string `mapstructure:"service_version"`
}
