package yeeUtil

import (
	"github.com/cqlsy/yeelib/yeecrypto"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

//生成token值
func RandString(length int, userId string) string {
	rand.Seed(time.Now().UnixNano())
	rs := make([]string, 0)
	//  生成随机字符串
	for start := 0; start < length; start++ {
		t := rand.Intn(3)
		if t == 0 {
			rs = append(rs, strconv.Itoa(rand.Intn(10)))
		} else if t == 1 {
			rs = append(rs, string(rune(rand.Intn(26)+65)))
		} else {
			rs = append(rs, string(rune(rand.Intn(26)+97)))
		}
	}
	// 加上时间戳以及唯一的参数,保证获取的哈希数据的唯一性
	rs = append(append(rs, userId), time.Now().String())
	return yeecrypto.Sha256Hex([]byte(strings.Join(rs, "")))
}
