package nacos

import (
	"context"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
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
			logger.SErrorF(ctx, "Nacos SearchConfig, search config failed. searchConfigParam: %+v, err: %v", searchConfigParam, err)
			return nil, err
		}

		logger.SInfoF(ctx, "Nacos SearchConfig, search config success. searchConfigParam: %+v, total count: %d, page number: %d, pages available: %d, data count:%d", searchConfigParam, searchResult.TotalCount, searchResult.PageNumber, searchResult.PagesAvailable, len(searchResult.PageItems))

		result = append(result, searchResult.PageItems...)
		pageNo = pageNo + 1

		if searchResult.PageNumber >= searchResult.PagesAvailable {
			break
		}
	}

	logger.SInfoF(ctx, "Nacos SearchConfig, search config success. group: %s, config count:%d", group, len(result))
	return result, nil
}
