package middle

import (
	"github.com/gin-gonic/gin"
	"github.com/leo/leolib/leocrypto/rsa"
	"github.com/leo/leolib/leofile"
	"github.com/leo/leolib/leolog"
	"github.com/leo/leolib/leoutil"
	"github.com/leo/leolib/webGin_example/server/base"
	"time"
)

var RsaPath = "./data/crypto/"
var PublicKey string

// change rsa every 24 hours
func Init(runMode string) {
	var work = func() {
		err := leofile.Mkdir(RsaPath)
		if err != nil {
			leolog.LogErrorDefault(err.Error())
			return
		}
		err = rsa.GenRsaKey(RsaPath, RsaPath, 1024*4)
		if err != nil {
			leolog.LogErrorDefault(err.Error())
		} else {
			PublicKey, _ = rsa.GetPubKeyString(RsaPath + "public.pem")
		}
	}
	if runMode == "dev" {
		if !leofile.FileExists(RsaPath + "public.pem") {
			work()
		} else {
			PublicKey, _ = rsa.GetPubKeyString(RsaPath + "public.pem")
		}
		//} else if runMode == "pro" {
	} else {
		leoutil.StartTickerTask(time.Hour*24, true, work)
	}
}

func CheckLogin() func(*gin.Context) {
	return func(c *gin.Context) {
		//c.Next() // 在它之后的代码将会在执行完其他的请求之后再执行
		//c.JSON()
		token := c.GetHeader("token")
		if token == "" {
			c.JSON(200, base.LoginInvalid())
			c.Abort() // end this request
		}
		c.Set("asekey", "")
	}
}
