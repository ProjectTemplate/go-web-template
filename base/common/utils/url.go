package utils

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"net/url"
	"reflect"
)

func MapToUrlValues(data map[string]interface{}) url.Values {
	values := url.Values{}
	for k, v := range data {
		if actualV, ok := v.(string); ok {
			values.Add(k, actualV)
			continue
		}

		vValue := reflect.ValueOf(v)
		if vValue.Kind() == reflect.Ptr {
			vValue = vValue.Elem()
		}

		if vValue.Kind() == reflect.Slice {
			if vValue.IsNil() {
				continue
			}
			if vValue.Len() < 0 {
				continue
			}

			length := vValue.Len()
			for i := 0; i < length; i++ {
				element := vValue.Index(i)
				values.Add(fmt.Sprintf("%s[]", k), cast.ToString(element.Interface()))
			}

			continue
		}

		toStringE, err := cast.ToStringE(v)
		if err == nil {
			values.Add(k, toStringE)
			continue
		}

		marshal, err := json.Marshal(v)
		if err == nil {
			values.Add(k, string(marshal))
			continue
		}
	}
	return values
}
