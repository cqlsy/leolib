/**
 * Created by angelina-zf on 17/2/27.
 */

// leocrypto
// 用于hash加密的包
// 依赖： "golang.org/x/leocrypto/bcrypt"
package leocrypto

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

// Sha256Hex
func Sha256Hex(data []byte) string {
	out := sha256.Sum256(data)
	return hex.EncodeToString(out[:])
}

// Sha512Hex
func Sha512Hex(data []byte) string {
	out := sha512.Sum512(data)
	return hex.EncodeToString(out[:])
}

// Md5Hex 小写hex
func Md5Hex(data []byte) string {
	hash := md5.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}
