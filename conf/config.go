package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Info struct {
	Db  db
	Web web
}

type web struct {
	LogPath      string
	RunMode      string
	Ip           string
	Port         int
	SaveFilePath string
}

type db struct {
	Host     string
	Db       string
	Port     string
	User     string
	Password string
	Protocol string
}

var Conf *Info // 全局实例，全局可访问。

// conf 需要传入地址
func ParseConf(confPath string) {
	Conf = new(Info)
	data, err := ioutil.ReadFile(confPath)
	if err != nil {
		panic(fmt.Sprintf("sync Info file read error: %s", err.Error()))
	}
	err = json.Unmarshal(data, &Conf)
	if err != nil {
		panic(fmt.Sprintf("sync Info format error: %s", err.Error()))
	}
}
