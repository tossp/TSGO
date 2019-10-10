package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/pbkdf2"
)

func HashPasswd(password, salt string) string {
	dk := argon2.IDKey(Sha512([]byte(password)), Sha512([]byte(salt)), 1, 32*1024, 4, 32)
	pwd := base64.StdEncoding.EncodeToString(dk[:])
	return pwd
}

func Sha512(input []byte) []byte {
	hmac512 := hmac.New(sha512.New, []byte("TossP.com"))
	hmac512.Write(input)
	bs := hmac512.Sum(nil)
	return bs[:]
}

func Hash32(password, salt []byte) []byte {
	return pbkdf2.Key(password, salt, 100, 32, sha256.New)[:]
}

func HashKey(input []byte, keylen int) (key []byte) {
	return pbkdf2.Key(input, []byte("TossP.com"), 1024, keylen, sha256.New)[:]
}
