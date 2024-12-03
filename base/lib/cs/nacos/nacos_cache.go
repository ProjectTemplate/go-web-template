package nacos

import (
	"context"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"go-web-template/base/lib/logger"
	"go.uber.org/zap"
	"strings"
	"sync"
)

// ConfigCache 配置数据缓存
//
// 使用 nacos 作为配置中心，缓存配置数据
//
// 通过 Listener 监听配置变化，实时更新
type ConfigCache struct {
	configClient config_client.IConfigClient

	cacheDataMap sync.Map //缓存数据 map[string]*cacheData
}

type cacheData struct {
	state             cacheState  //缓存状态
	data              interface{} //缓存数据
	unmarshalDataFunc UnmarshalDataFunc
}

func NewConfigCache(configClient config_client.IConfigClient) *ConfigCache {
	result := &ConfigCache{}
	result.configClient = configClient
	return result
}

// InitConfig  初始化配置
func (n *ConfigCache) InitConfig(ctx context.Context, group, dataId string, unmarshalFunc UnmarshalDataFunc) {
	logger.Info(ctx, "InitConfig init nacos config.", zap.String("group", group), zap.String("dataId", dataId), zap.Any("unmarshalFunc", unmarshalFunc))

	n.listenConfig(ctx, group, dataId)

	//设置缓存和默认值
	data := &cacheData{
		state:             cacheState{false},
		unmarshalDataFunc: unmarshalFunc,
		data:              nil,
	}
	n.cacheDataMap.Store(dataUid(group, dataId), data)
	configData, err := n.configClient.GetConfig(vo.ConfigParam{DataId: dataId, Group: group})
	if err != nil {
		logger.Warn(ctx, "InitConfig get nacos config failed.", zap.String("group", group), zap.String("dataId", dataId), zap.String("configData", configData))
		return
	}

	unmarshalResult, err := unmarshalFunc(configData)
	if err != nil {
		logger.Warn(ctx, "InitConfig write config data failed.",
			zap.String("group", group), zap.String("dataId", dataId), zap.String("configData", configData))
	} else {
		data.data = unmarshalResult
		data.state.enable()
	}
}

// GetConfig 返回配置数据，返回之后用类型断言断言类型
// 默认情况改接口不会返回空值
func (n *ConfigCache) GetConfig(ctx context.Context, group, dataId string) interface{} {
	value, ok := n.cacheDataMap.Load(dataUid(group, dataId))
	if !ok {
		logger.Info(ctx, "GetConfig config cache not found.", zap.String("group", group), zap.String("dataId", dataId))
		return nil
	}
	cData := value.(*cacheData)
	if cData.state.active() {
		return cData.data
	}
	return nil
}

// listenConfig 监听Nacos配置变化，同时设置缓存数据
func (n *ConfigCache) listenConfig(ctx context.Context, group, dataId string) {
	err := n.configClient.ListenConfig(vo.ConfigParam{
		DataId:   dataId,
		Group:    group,
		OnChange: n.updateConfig(ctx),
	})
	if err != nil {
		logger.Warn(ctx, "listenConfig listen config failed.", zap.String("group", group), zap.String("dataId", dataId))
	}
}

func (n *ConfigCache) updateConfig(ctx context.Context) func(namespace, group, dataId, data string) {
	return func(namespace, group, dataId, data string) {
		logger.Info(ctx, "updateConfig listen config change.",
			zap.String("namespace", namespace), zap.String("group", group), zap.String("dataId", dataId), zap.String("data", data))

		cDataInterface, ok := n.cacheDataMap.Load(dataUid(group, dataId))
		if !ok {
			logger.Info(ctx, "updateConfig config cache not found.",
				zap.String("namespace", namespace), zap.String("group", group), zap.String("dataId", dataId), zap.String("data", data))
			return
		}

		cData := cDataInterface.(*cacheData)
		unmarshalData, err := cData.unmarshalDataFunc(data)
		if err != nil {
			cData.state.disable()
			logger.Warn(ctx, "updateConfig listen config unmarshal data filed.",
				zap.String("namespace", namespace), zap.String("group", group), zap.String("dataId", dataId), zap.String("data", data))
			return
		}
		cData.data = unmarshalData
		cData.state.enable()
	}
}

func dataUid(groupId string, dataId string) string {
	return strings.Join([]string{groupId, dataId}, "_")
}
