package nacos

import (
	"context"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"go.uber.org/zap"

	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"
)

const (
	// dir nacos目录
	dir = "/tmp/nacos"
	// cacheDir 缓存目录
	cacheDir = dir + "/cache"
	// LogDir 日志目录，nacos 和 基础数据sdk放在同一个目录
	LogDir = dir + "/log"

	// LogMaxSize 日志文件大小
	LogMaxSize = 100
	// LogMaxAge 日志文件最大保存时间
	LogMaxAge = 15

	// searchPageSize 搜素配置的每页大小
	searchPageSize = 10

	// ClientTypeConfig 配置客户端
	ClientTypeConfig ClientType = "config"
	// ClientTypeNaming 命名客户端
	ClientTypeNaming ClientType = "naming"
	// ClientTypeAll 所有客户端
	ClientTypeAll ClientType = "all"
)

var cli *client

type ClientType string

// Client nacos客户端
type client struct {
	configClients map[string]config_client.IConfigClient
	namingClients map[string]naming_client.INamingClient
}

// GetConfigClient 根据名字获取配置客户端 [config_client.IConfigClient]，如果不存在则panic
func GetConfigClient(ctx context.Context, name string) config_client.IConfigClient {
	configClient := cli.configClients[name]
	if configClient == nil {
		logger.Error(ctx, "GetConfigClient, config client not found", zap.String("name", name))
		panic("GetConfigClient, config client not found, name: " + name)
	}
	return configClient
}

// GetNamingClient 根据名字获取命名客户端 [naming_client.INamingClient]，如果不存在则panic
func GetNamingClient(ctx context.Context, name string) naming_client.INamingClient {
	namingClient := cli.namingClients[name]
	if namingClient == nil {
		logger.Error(ctx, "GetNamingClient, naming client not found", zap.String("name", name))
		panic("GetNamingClient, naming client not found, name: " + name)
	}
	return namingClient
}

// Init 初始化nacos客户端
func Init(ctx context.Context, nacosConfigs map[string]config.Nacos) {
	logger.Info(ctx, "Init Nacos start.", zap.Any("nacosConfigs", nacosConfigs))

	cli = &client{
		configClients: make(map[string]config_client.IConfigClient),
		namingClients: make(map[string]naming_client.INamingClient),
	}

	for name, nacosConfig := range nacosConfigs {
		clientType := nacosConfig.ClientType
		if clientType != string(ClientTypeAll) && clientType != string(ClientTypeConfig) && clientType != string(ClientTypeNaming) {
			logger.Error(ctx, "Init Nacos, client type error, client type must be [config naming all].", zap.String("name", name), zap.Any("config", nacosConfig))
			panic("Init Nacos failed, client type error, client type must be [config naming all], type: " + clientType)
		}

		if clientType == string(ClientTypeAll) || clientType == string(ClientTypeConfig) {
			configClient, err := newConfigClient(ctx, nacosConfig)
			if err != nil {
				logger.Error(ctx, "Nacos Init, new config client failed.", zap.String("name", name), zap.Any("config", nacosConfig), zap.Error(err))
				panic("Init Nacos failed, err:" + err.Error())
			}
			cli.configClients[name] = configClient
		}

		if clientType == string(ClientTypeAll) || clientType == string(ClientTypeNaming) {
			namingClient, err := newNamingClient(ctx, nacosConfig)
			if err != nil {
				logger.Error(ctx, "Init Nacos, new config client failed.", zap.String("name", name), zap.Any("config", nacosConfig), zap.Error(err))
				panic("Init Nacos failed, err:" + err.Error())
			}
			cli.namingClients[name] = namingClient
		}

		logger.Info(ctx, "Init Nacos success.", zap.String("name", name), zap.Any("config", nacosConfig))
	}
}

// newConfigClient 创建配置客户端
func newConfigClient(ctx context.Context, nacosConfig config.Nacos) (config_client.IConfigClient, error) {
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  initClientConfig(nacosConfig.Namespace),
			ServerConfigs: initServerConfigs(nacosConfig.Servers),
		},
	)
	if err != nil {
		logger.Error(ctx, "Nacos NewConfigClient, new config client failed.", zap.String("nameSpaceId", nacosConfig.Namespace), zap.Any("servers", nacosConfig.Servers), zap.Error(err))
		return configClient, err
	}

	return configClient, nil
}

func newNamingClient(ctx context.Context, nacosConfig config.Nacos) (naming_client.INamingClient, error) {
	configClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  initClientConfig(nacosConfig.Namespace),
			ServerConfigs: initServerConfigs(nacosConfig.Servers),
		},
	)

	if err != nil {
		logger.Error(ctx, "Nacos NewNamingClient, new naming client failed.", zap.String("nameSpaceId", nacosConfig.Namespace), zap.Any("servers", nacosConfig.Servers), zap.Error(err))
		return configClient, err
	}

	return configClient, nil
}

func initClientConfig(nameSpaceId string) *constant.ClientConfig {
	clientConfig := constant.ClientConfig{
		NamespaceId: nameSpaceId,

		TimeoutMs: 60000,

		NotLoadCacheAtStart: true,

		CacheDir: cacheDir,

		LogLevel:       "info",
		LogDir:         LogDir,
		AppendToStdout: false,
		LogRollingConfig: &constant.ClientLogRollingConfig{
			MaxSize:   LogMaxSize,
			MaxAge:    LogMaxAge,
			LocalTime: true,
		},
	}

	return &clientConfig
}

// initServerConfigs 创建服务配置
func initServerConfigs(servers []config.NacosServerConf) []constant.ServerConfig {
	serverConfigs := make([]constant.ServerConfig, 0, len(servers))

	for _, server := range servers {
		serverConfigs = append(serverConfigs, constant.ServerConfig{
			Scheme:      "http",
			IpAddr:      server.IpAddr,
			Port:        server.Port,
			ContextPath: "/nacos",
		})
	}

	return serverConfigs
}
