package leoutil

import (
	"errors"
	"fmt"
	"strconv"
)

func GetString(data interface{}) string {
	if data == nil {
		return ""
	}
	switch data.(type) {
	case string:
		return data.(string)
	}
	return fmt.Sprintf("%v", data)
}

func GetInt64WithDefault(data interface{}, def int64) int64 {
	result, err := GetInt64(data)
	if err != nil {
		return def
	}
	return result
}

func GetInt64(data interface{}) (int64, error) {
	if data ==nil {
		return 0,errors.New("the data is nil")
	}
	switch data.(type) {
	case int:
		return int64(data.(int)), nil
	case int64:
		return data.(int64), nil
	case int8:
		return int64(data.(int8)), nil
	case int16:
		return int64(data.(int16)), nil
	case int32:
		return int64(data.(int32)), nil
	case string:
		return strconv.ParseInt(data.(string), 10, 64)
	case float32:
		return int64(data.(float32)), nil
	case float64:
		return int64(data.(float64)), nil
	}
	return 0, errors.New("error type ")
}
