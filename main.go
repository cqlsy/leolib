package main

import (
	"github.com/cqlsy/leolib/args"
	"github.com/cqlsy/leolib/leoconfig"
	"github.com/cqlsy/leolib/leodb"
	fileUtil "github.com/cqlsy/leolib/leofile/action"
	"github.com/cqlsy/leolib/leolog"
	"github.com/cqlsy/leolib/leowebgin"
	webEcho "github.com/cqlsy/leolib/webEcho_example"
	webGin "github.com/cqlsy/leolib/webGin_example"
	"github.com/cqlsy/leolib/webGin_example/server/base"
)

var (
	conf       *leoconfig.Info
	configPath = "./config.json"
)

func main() {
	if args.GetArgs(func(key string, value string) (b bool, b2 bool) {
		switch key {
		}
		return
	}) {
		return
	}

	// conf
	conf = new(leoconfig.Info)
	leoconfig.ParseConf(configPath, &conf)
	// log
	leolog.MustInitLog(conf.Web.LogPath, conf.Web.RunMode)
	// file path for upload and download
	fileUtil.Init(conf.Web.SaveFilePath)
	// database
	database := leodb.InitDataBase(configPath)

	// service init
	base.InitPicDoMain(conf.Web.PicDoMain)
	// web service
	startGinService(database)
}

func startGinService(db *leodb.Db) {
	ginInstance := leowebgin.New(conf.Web.RunMode)
	ginInstance.AddStaticPath("/data", "./data")
	if conf.Web.RunMode == "dev" {
		ginInstance.AddStaticPath("/app", "./static")
	} else {
		ginInstance.AddStaticPackrPath("/app", "./static", "./static")
	}
	webGin.Init(ginInstance, db, conf.Web.RunMode)
	//ginInstance.AddStaticPath("/app", "./static")
	ginInstance.StartListen(conf.Web.Ip, conf.Web.Port)
}

func startEchoService() {
	webEcho.StartWeb()
}
