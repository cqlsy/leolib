/**
 * Created by angelina-zf on 17/2/25.
 */
package leolog

import (
	"testing"
	"time"
)

func TestLogDebug(t *testing.T) {
	//MustInitLog("./logs/", "pro")
	/*Log("debug", LogFields{
		//"event": "a",
		//"topic": "b",
		//"key":   "C",
	},"test")*/
	Print("hello","dsahlkhda")
}

func TestLogError(t *testing.T) {
	//LogError("error", LogFields{
	//	"event": "a",
	//	"topic": "b",
	//	"key":   "C",
	//})
}

func TestLogInfo(t *testing.T) {
	//LogInfo("info", LogFields{
	//	"event": "a",
	//	"topic": "b",
	//	"key":   "C",
	//})
}

func TestLogFile(t *testing.T) {
	MustInitLog("logs/", "pro")
	//go func() {
	//	for i := 1; i < 10000; i++ {
	//		LogDefault("test", "test")
	//	}
	//}()
	//LogDefault("TEST1", "test1")
	time.Sleep(time.Second * 10)
}

func TestFile(t *testing.T) {
	//if leofile.FileExists("../data/log") { // 检查是否有当前的目录
	//	files, _ := ioutil.ReadDir("../data/log")
	//	for _, f := range files {
	//		println(f.Name())
	//		t, err := time.Parse("2006-01-02", f.Name())
	//		if err != nil {
	//			println(err.Error())
	//		} else {
	//			now := time.Now().Unix()
	//			if (now - t.Unix()) > 60*60*24*5 {
	//				err := os.RemoveAll("../data/log" + "/" + f.Name())
	//				if err != nil {
	//					LogDefault("Delete ExpiredLog Err: "+err.Error(), config.BlockLatest)
	//				}
	//			}
	//		}
	//	}
	//}
}
