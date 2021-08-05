package yeeGin

import (
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	config "lib/conf"
	"lib/yeelog"
	"net/http"
	"time"
)

type WebGin struct {
	Gin *gin.Engine
}

// gin的取值方式，记录
//			c *gin.Context
// 1: filePath := c.Query("filePath")   	// 获取 Query（拼接到Url后面的参数）
// 2: filePath = c.PostForm("filePath") 	// Post 的请求参数
// 3: data := make(map[string]interface{})  // body中的参数
//    c.ShouldBind(&data)
func New() *WebGin {
	//gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard // 不让它自身打印,使用我们自己的打印
	web := new(WebGin)
	web.Gin = gin.Default()
	web.Gin.Use(cors(), gzip.Gzip(gzip.DefaultCompression))
	if config.Conf.Web.RunMode != "pro" {
		web.Gin.Use(middleLog())
	}
	return web
}

func cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token ,X-Token")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
		context.Next()
	}
}

// 开启监听 ip:0.0.0.0  port:8080
func (w *WebGin) StartListen(ip interface{}, port interface{}) {
	str := fmt.Sprintf("%v:%v", ip, port)
	yeelog.Print("Web Listener On:" + str)
	err := w.Gin.Run(str)
	if err != nil {
		panic(fmt.Sprintf("Listener %s error: %s", str, err.Error()))
	}
}

// packr.New("static", "./static") // 资源打包到二进制文件时，静态路由设置
//  http.Dir("./static") // 普通路由设置
func (w *WebGin) AddStaticPath(path string, fs http.FileSystem) {
	w.Gin.StaticFS(path, fs)
}

func middleLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.Request.Host
		// 日志格式
		yeelog.LogInfoDefault(fmt.Sprintf("| code: %3d | time: %13v | clientIP: %15s | reqMethod: %s | reqUri: %s |",
			statusCode, latencyTime, clientIP, reqMethod, reqUri))
	}
}

// 两种传值的方式全部使用到了
func GetReqParams(c *gin.Context) string {
	params := c.Query("filePath") // 获取 Query（拼接到Url后面的参数）
	if params == "" {
		params = c.PostForm("filePath") // Post 的请求参数
	}
	return params
}

// 从Body中取值
func GetReqParamsFromBody(c *gin.Context) map[string]interface{} {
	data := make(map[string]interface{})
	_ = c.ShouldBind(&data)
	return data
}
