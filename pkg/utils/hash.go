package utils

import (
	"github.com/tossp/tsgo/pkg/utils/crypto"
)

func HashPasswd(password string) string {
	pwd, _ := GenerateHash(password, "argon2")
	return pwd
}
func ComparePasswd(password, encodedHash string) (match bool, err error) {
	return CompareHash(password, encodedHash, "argon2")
}

func GenerateHash(input, hasherName string) (string, error) {
	hasher := crypto.GetHasher(hasherName)
	return hasher.Generate(input)
}
func CompareHash(input, encodedHash, hasherName string) (bool, error) {
	hasher := crypto.GetHasher(hasherName)
	return hasher.Compare(input, encodedHash)
}
