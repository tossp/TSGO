package utils

import (
	"github.com/tossp/tsgo/pkg/utils/crypto"
)

//GetRandomString 通过指定字符集，生成随机字符串
func GetRandomString(n int, alphabets ...byte) string {
	return crypto.GetRandomString(n, alphabets...)
}
