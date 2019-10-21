package crypto

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/pbkdf2"
)

func Sha1(input []byte) []byte {
	hmac1 := hmac.New(sha1.New, []byte("TossP.com"))
	hmac1.Write(input)
	bs := hmac1.Sum(nil)
	return bs[:]
}
func Sha256(input []byte) []byte {
	hmac256 := hmac.New(sha256.New, []byte("TossP.com"))
	hmac256.Write(input)
	bs := hmac256.Sum(nil)
	return bs[:]
}
func Sha512(input []byte) []byte {
	hmac512 := hmac.New(sha512.New, []byte("TossP.com"))
	hmac512.Write(input)
	bs := hmac512.Sum(nil)
	return bs[:]
}

func Hash32(password, salt []byte) []byte {
	return pbkdf2.Key(password, Sha256(salt), 100, 32, sha256.New)[:]
}

func HashKey(input []byte, keylen int) (key []byte) {
	return pbkdf2.Key(input, Sha1([]byte("TossP.com")), 1024, keylen, sha256.New)[:]
}
func HashSha(input, salt []byte, keylen int) (key []byte) {
	return pbkdf2.Key(input, Sha512(salt), 1024, keylen, sha512.New)[:]
}
func HashArgon(input, salt []byte, keylen uint32) (key []byte) {
	return argon2.Key(input, Sha512(salt), 4, 32*1024, 3, keylen)
}
