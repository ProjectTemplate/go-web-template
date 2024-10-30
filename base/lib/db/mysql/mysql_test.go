package mysql

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"
	"testing"
)

func TestGorm(t *testing.T) {
	configStruct := &config.Configs{}
	config.Init("./data/config.toml", configStruct)
	logger.Init(configStruct.LoggerConfig)

	dbMap, err := Init(configStruct.DB)

	assert.Nil(t, err)
	assert.NotNil(t, dbMap)
	assert.NotNil(t, dbMap["test"])

	gormDB := dbMap["test"]
	count := new(int64)
	err = gormDB.Table("person").Count(count).Error
	assert.Nil(t, err)
	fmt.Println("count: ", *count)
}
