package mysql

import (
	"context"
	"fmt"
	"testing"

	"go-web-template/base/lib/config"
	"go-web-template/base/lib/logger"

	"github.com/stretchr/testify/assert"
)

func TestCount(t *testing.T) {
	configStruct := &config.Configs{}
	config.Init("./data/config.toml", configStruct)
	logger.Init("TestGorm", configStruct.LoggerConfig)

	background := context.Background()

	Init(background, configStruct.Mysql)

	db := GetDB(background, "test")
	assert.NotNil(t, db)

	count := new(int64)
	err := db.Table("person").Count(count).Error
	assert.Nil(t, err)
	fmt.Println("count: ", *count)
}

func TestGormError(t *testing.T) {
	configStruct := &config.Configs{}
	config.Init("./data/config.toml", configStruct)
	logger.Init("TestGormError", configStruct.LoggerConfig)

	background := context.Background()
	Init(background, configStruct.Mysql)
	db := GetDB(background, "test")

	assert.NotNil(t, db)

	count := new(int64)
	err := db.Exec("show tables").Count(count).Error
	assert.NotNil(t, err)
}
