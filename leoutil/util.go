package leoutil

import (
	"encoding/json"
	"errors"
	"github.com/cqlsy/leolib/leocrypto"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
)

//生成token值
func RandString(length int, userId string) string {
	rand.Seed(time.Now().UnixNano())
	rs := make([]string, 0)
	//  生成随机字符串
	for start := 0; start < length; start++ {
		t := rand.Intn(3)
		if t == 0 {
			rs = append(rs, strconv.Itoa(rand.Intn(10)))
		} else if t == 1 {
			rs = append(rs, string(rune(rand.Intn(26)+65)))
		} else {
			rs = append(rs, string(rune(rand.Intn(26)+97)))
		}
	}
	// 加上时间戳以及唯一的参数,保证获取的哈希数据的唯一性
	rs = append(append(rs, userId), time.Now().String())
	// 返回hash256的字符串
	return leocrypto.Sha256Hex([]byte(strings.Join(rs, "")))
}

func Struct2Map(obj interface{}) (map[string]interface{}, error) {
	var data map[string]interface{}
	jsonData, err := json.Marshal(obj)
	if err == nil {
		err = json.Unmarshal(jsonData, &data)
	}
	return data, err
}

func GetStructValueByJsonId(data interface{}, key string) (interface{}, error) {
	resultsVal := reflect.ValueOf(data)
	if resultsVal.Kind() == reflect.Ptr {
		return nil, errors.New("can't user Ptr to get Data")
	}
	types := reflect.TypeOf(data)
	name := ""
	for index := 0; index < types.NumField(); index++ {
		ss := strings.Split(types.Field(index).Tag.Get("json"), ",")
		for _, item := range ss {
			if item == key {
				name = types.Field(index).Name
				break
			}
		}
	}
	if name == "" {
		return nil, errors.New("no Key with the key:" + key)
	}
	ss := reflect.ValueOf(data)
	return ss.FieldByName(name).Interface(), nil
}

func Struct2MapD(obj interface{}) map[string]interface{} {
	data := make(map[string]interface{})
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func Map2Struct(src interface{}, dst interface{}) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, dst)
	return err
}

func ObjectToJSONStr(data interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(b)
}
