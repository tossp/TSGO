package utils

import "crypto/rand"

const (
	alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

//GetRandomString 通过指定字符集，生成随机字符串
func GetRandomString(n int, alphabets ...byte) string {
	var bytes = make([]byte, n)
	_, _ = rand.Read(bytes)
	if len(alphabets) == 0 {
		alphabets = []byte(alphanum)
	}
	for i, b := range bytes {
		bytes[i] = alphabets[b%byte(len(alphabets))]
	}
	return string(bytes)
}
