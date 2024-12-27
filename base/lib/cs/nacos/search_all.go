package nacos

import (
	"context"

	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"go.uber.org/zap"

	"go-web-template/base/lib/logger"
)

// SearchConfigInGroup 搜索足内的所有配置
func SearchConfigInGroup(ctx context.Context, group string, configClient config_client.IConfigClient) ([]model.ConfigItem, error) {
	result := make([]model.ConfigItem, 0)

	pageNo := 1
	for {

		searchConfigParam := vo.SearchConfigParam{
			Search:   "blur",
			Group:    group,
			DataId:   "",
			PageNo:   pageNo,
			PageSize: searchPageSize,
		}

		searchResult, err := configClient.SearchConfig(searchConfigParam)

		if err != nil {
			logger.Error(ctx, "Nacos SearchConfig, search config failed.", zap.Any("searchConfigParam", searchConfigParam), zap.Error(err))
			return nil, err
		}

		logger.Info(ctx, "Nacos SearchConfig, search config success.",
			zap.Any("searchConfigParam", searchConfigParam),
			zap.Int("total count", searchResult.TotalCount),
			zap.Int("page number", searchResult.PageNumber),
			zap.Int("pages available", searchResult.PagesAvailable),
			zap.Int("data count", len(searchResult.PageItems)))

		result = append(result, searchResult.PageItems...)
		pageNo = pageNo + 1

		if searchResult.PageNumber >= searchResult.PagesAvailable {
			break
		}
	}

	logger.Info(ctx, "Nacos SearchConfig", zap.String("group", group), zap.Int("configCount", len(result)))
	return result, nil
}
