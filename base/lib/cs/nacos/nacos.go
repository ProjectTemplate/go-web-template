package nacos

import (
	"context"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
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

type ClientType string

// Client nacos客户端
type Client struct {
	configClients map[string]config_client.IConfigClient
	namingClients map[string]naming_client.INamingClient
}

// Init 初始化nacos客户端
func Init(ctx context.Context, nacosConfigs map[string]config.Nacos) Client {
	logger.SInfoF(ctx, "Nacos Init start, nacos configs: %+v", nacosConfigs)

	client := Client{
		configClients: make(map[string]config_client.IConfigClient),
		namingClients: make(map[string]naming_client.INamingClient),
	}

	for name, nacosConfig := range nacosConfigs {
		clientType := nacosConfig.ClientType
		if clientType != string(ClientTypeAll) && clientType != string(ClientTypeConfig) && clientType != string(ClientTypeNaming) {
			logger.SErrorF(ctx, "Nacos Init, client type error, client type must be [config naming all]. name: %s, config:%+v", name, nacosConfig)
			panic("Init Nacos failed, client type error, client type must be [config naming all], type: " + clientType)
		}

		if clientType == string(ClientTypeAll) || clientType == string(ClientTypeConfig) {
			configClient, err := newConfigClient(ctx, nacosConfig)
			if err != nil {
				logger.SErrorF(ctx, "Nacos Init, new config client failed. name: %s, config:%+v err: %v", name, nacosConfig, err)
				panic("Init Nacos failed, err:" + err.Error())
			}
			client.configClients[name] = configClient
		}

		if clientType == string(ClientTypeAll) || clientType == string(ClientTypeNaming) {
			namingClient, err := newNamingClient(ctx, nacosConfig)
			if err != nil {
				logger.SErrorF(ctx, "Nacos Init, new config client failed. name: %s, config:%+v err: %v", name, nacosConfig, err)
				panic("Init Nacos failed, err:" + err.Error())
			}
			client.namingClients[name] = namingClient
		}

		logger.SInfoF(ctx, "Nacos Init success, name: %s, config:%+v", name, nacosConfig)
	}

	return client
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
		logger.SErrorF(ctx, "Nacos NewConfigClient, new config client failed. nameSpaceId: %s, servers: %+v, err: %v", nacosConfig.Namespace, nacosConfig.Servers, err)
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
		logger.SErrorF(ctx, "Nacos NewNamingClient, new naming client failed. nameSpaceId: %s, servers: %+v, err: %v", nacosConfig.Namespace, nacosConfig.Servers, err)
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
