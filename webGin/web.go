package webGin

import (
	config "github.com/lsy/yeelib/conf"
	"github.com/lsy/yeelib/webGin/server"
	"github.com/lsy/yeelib/yeeDb"
	"github.com/lsy/yeelib/yeeGin"
	"github.com/lsy/yeelib/yeefile"
)

// 我们使用这个来保存web服务需要的实例
var MongoDb *yeeDb.Db     // 数据库查询
var Engine *yeeGin.WebGin // 设置路由

func Init(engine *yeeGin.WebGin, db *yeeDb.Db) {
	MongoDb = db
	Engine = engine
	InitRouter()
	yeefile.Init(config.Conf.Web.SaveFilePath)
}

// 初始化路由信息
func InitRouter() {
	//Engine.Gin.Group()
	gin := Engine.Gin
	gin.POST("/Test", server.Test{}.NetTest)
	group := gin.Group("/group")
	group.POST("/Test", server.Test{}.NetTest)

	file := server.File{}
	gin.POST("/uploadFile", file.UploadFile)
	gin.POST("/downloadFile", file.DownLoadFile)

}
