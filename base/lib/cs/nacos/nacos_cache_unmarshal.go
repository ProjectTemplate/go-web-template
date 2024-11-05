package nacos

import (
	"encoding/json"
	"errors"
	"github.com/spf13/cast"
)

type UnmarshalDataFunc func(data string) (interface{}, error)

func UnmarshalToMap(data string) (interface{}, error) {
	result := make(map[string]string)
	err := json.Unmarshal([]byte(data), &result)
	return result, err
}

func UnmarshalToBool(data string) (interface{}, error) {
	if data == "1" {
		return true, nil
	}
	if data == "0" {
		return false, nil
	}
	return nil, errors.New("unmarshal to bool failed. data: " + data)
}

func UnmarshalToNumber(data string) (interface{}, error) {
	return cast.ToIntE(data)
}
