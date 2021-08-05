package yeeDb

import (
	config "lib/conf"
	"lib/yeeDb/mog"
	"strings"
)

// 保存我们的数据实例
type Db struct {
	MogDb *mog.MongoDb
}

// 程序在这里初始化链接数据库
func InitDataBase() *Db {
	if config.Conf == nil {
		panic("请完成数据的相关配置")
	}
	db := new(Db)
	if strings.ToUpper(config.Conf.Db.Protocol) == "MONGODB" {
		mon, err := mog.Collect()
		if err != nil {
			panic("数据库链接失败： " + err.Error())
		}
		db.MogDb = mon
	} else if strings.ToUpper(config.Conf.Db.Protocol) == "MYSQL" {

	}
	return db
}
