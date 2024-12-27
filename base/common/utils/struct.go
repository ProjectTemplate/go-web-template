package utils

import (
	"github.com/fatih/structs"

	"go-web-template/base/common/constant"
)

func StructToMap(data interface{}) map[string]interface{} {
	s := structs.New(data)
	s.TagName = constant.TagNameJson
	return s.Map()
}
