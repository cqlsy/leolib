package webGin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/cqlsy/leolib/leodb"
	"github.com/cqlsy/leolib/leolog"
	"github.com/cqlsy/leolib/leowebgin"
	"github.com/cqlsy/leolib/webGin_example/model"
	"github.com/cqlsy/leolib/webGin_example/server"
	"github.com/cqlsy/leolib/webGin_example/server/middle"
	"net/http"
)

var Db *leodb.Db // should Save it on here
var Engine *leowebgin.WebGin

func Init(engine *leowebgin.WebGin, db *leodb.Db, runMode string) {
	Db = db
	Engine = engine
	model.Init(Db)
	middle.Init(runMode)
	initRouter()
	initSocket()
}

// manager save to use
func initSocket() {
	manager := leowebgin.NewManager(
		nil,
		func(client *leowebgin.Client, msg []byte) {
			leolog.LogInfoDefault(string(msg))
			client.SendMessage([]byte(fmt.Sprintf("service callback：%s", msg)))
		},
	)
	// manager.Clients
	// socket
	Engine.AddSocketClient("/socket", manager, nil)
}

func initRouter() {
	//Engine.Gin.Group()
	g := Engine.Gin
	//group := g.Group("/group")
	server.Album{}.Init(g)
	server.Photo{}.Init(g)

	file := server.File{}
	g.GET("/file/:fileName", file.DownLoadFile)
	g.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/app")
		c.Abort()
	})
	//g.POST("/uploadFile", file.UploadFile)
	//g.GET("/", func(c *gin.Context) {
	//	c.Redirect(http.StatusMovedPermanently, "http://www.baidu.com")
	//})
	//g.POST("/test", middle.CheckLogin(), server.User{}.Test())
	//g.POST("/getCryptoPubKey", server.User{}.GetRsaPubKey())

}
