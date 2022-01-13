package leoconfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// 	Example:
// 	conf := new(config.Info)
//	config.ParseConf(configPath, &conf)
func ParseConf(confPath string, conf interface{}) {
	data, err := ioutil.ReadFile(confPath)
	if err != nil {
		panic(fmt.Sprintf("sync Info file read error: %s", err.Error()))
	}
	err = json.Unmarshal(data, conf)
	if err != nil {
		panic(fmt.Sprintf("sync Info format error: %s", err.Error()))
	}
}
