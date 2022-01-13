/**
 * Created by angelina-zf on 17/2/27.
 */
package leotime

import (
	"testing"
	"time"
)

func TestDateFormat(t *testing.T) {
	println(DateFormat(time.Now(), "YYYY-MM-DD HH:mm:ss"))
}

func TestDay(t *testing.T) {
	var now = time.Now()
	println(now.Nanosecond())     // 只显示 秒之后的毫秒，微秒，纳秒 // 一般不用
	println(now.UnixNano() / 1e6) // 19位时间戳
	println(now.Unix())           // 10位时间戳
	println(now.String())
	// 时区设置,, 相对于的是 UTC的时区设置
	now = now.In(time.FixedZone("AAA", int((56 * time.Hour).Seconds())))
	println(now.String())
}
