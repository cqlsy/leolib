package yeeDb

import (
	config "lib/conf"
	"lib/yeeDb/mog"
	"lib/yeelog"
	"testing"
)

func TestSearch(t *testing.T) {
	config.Conf = new(config.Info)
	/**
	  "host": "192.168.56.101",
	    "port": "8666",
	    "db": "explorerDB",
	    "user": "explorer",
	    "password": "123456",
	    "protocol": "mongodb"
	*/
	config.Conf.Db.Host = "192.168.56.101"
	config.Conf.Db.Port = "8666"
	config.Conf.Db.Db = "explorerDB"
	config.Conf.Db.User = "explorer"
	config.Conf.Db.Password = "123456"
	config.Conf.Db.Protocol = "mongodb"
	// 手动设置配置文件
	database := InitDataBase()
	filter := mog.NewFilter()
	//filter.AddLikeCondition("symbol", "u")
	//filter.AddBiggerCondition("blockNumber", 426787, true)
	//filter.AddLessCondition("blockNumber", 426811, false)
	filter.AddExistsCondition("symbol", true)
	filter.AddNeCondition("symbol", "WT")

	//filter
	orFilter := mog.NewOrFilter()
	orA := mog.NewFilter()
	orA.AddEqualCondition("blockNumber", 426724)
	orFilter.AddCondition(orA)
	orA = mog.NewFilter()
	orA.AddEqualCondition("blockNumber", 426787)
	orFilter.AddCondition(orA)
	orA = mog.NewFilter()
	orA.AddLikeCondition("symbol", "UAT")
	orFilter.AddCondition(orA)
	filter.AddOrCondition(orFilter)
	_, err := database.MogDb.FindMany("Contract", []string{"symbol", "address", "blockNumber"}, map[string]int{}, nil, 0, 10)
	//if err == nil {
	//	yeelog.Print(data)
	//} else {
	//	yeelog.Print(err.Error())
	//}
	count, err := database.MogDb.Count("Contract", filter)
	if err == nil {
		yeelog.Print(count)
	} else {
		yeelog.Print(err.Error())
	}
	data := make(map[string]interface{})
	data["address"] = "0x222dasd222ss"
	data["blockNumber"] = 13123113
	data["symbol"] = "fasfasfasfasfasfas"

	//database.MogDb.UpdateOrInsertOne("Contract", data, "address")
	//database.MogDb.UpdateOrInsertOneWithKeys("Contract", data, []string{"address", "blockNumber"})
	datas := make([]map[string]interface{}, 0)
	datas = append(datas, data)
	data = make(map[string]interface{})
	data["address"] = "0x22222dasd"
	data["blockNumber"] = 13123112
	data["symbol"] = "fasfasfasfasfasfas"
	datas = append(datas, data)

	data = make(map[string]interface{})
	data["address"] = "0x222dasd"
	data["blockNumber"] = 13123112
	data["symbol"] = "mmm"
	datas = append(datas, data)
	//database.MogDb.UpdateOrInsertMany("Contract", datas, "address")
	//database.MogDb.UpdateOrInsertOneWithKeys("Contract", data, []string{"symbol", "blockNumber"})
	database.MogDb.UpdateOrInsertManyWithKeys("Contract", datas, []string{"address", "blockNumber"})

}
