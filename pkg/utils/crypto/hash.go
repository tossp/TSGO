package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"io"
	"os"

	"github.com/tjfoc/gmsm/sm3"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/pbkdf2"
)

const (
	defaultAlgo = "argon2"
)

var (
	implementations = make(map[string]hash)

	errUnknownHashImpl = errors.New("unknown tshash implementation")
	errBadHashFormat   = errors.New("invalid tshash format")
)

type hash interface {
	GetID() string
	String() string
	Generate(plainText string, param ...string) (encodedHash string, err error) //生成hash
	Compare(plainText, encodedHash string) (match bool, err error)              //比较hash
	Configure(param string) (hash, error)                                       //指定新的参数
	NeedsReHash(encodedHash string) (bool, error)                               //判断需要重新hash
}

func GetHasher(name string) hash {
	hasher, ok := implementations[name]
	if !ok {
		hasher = implementations[defaultAlgo]
	}
	return hasher
}

// GenerateRandomBytes Generate n number of random bytes
func GenerateRandomBytes(length int) ([]byte, error) {
	b := make([]byte, length)
	n, err := rand.Read(b)
	if err != nil || n != length {
		return nil, err
	}
	return b, nil
}

const (
	alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func GetRandomString(n int, alphabets ...byte) string {
	bytes, _ := GenerateRandomBytes(n)
	if len(alphabets) == 0 {
		alphabets = []byte(alphanum)
	}
	for i, b := range bytes {
		bytes[i] = alphabets[b%byte(len(alphabets))]
	}
	return string(bytes)
}

func Sha1(input []byte) []byte {
	hmac1 := hmac.New(sha1.New, []byte("TossP.com"))
	hmac1.Write(input)
	bs := hmac1.Sum(nil)[:]
	return bs
}
func Sha256(input []byte) []byte {
	hmac256 := hmac.New(sha256.New, []byte("TossP.com"))
	hmac256.Write(input)
	bs := hmac256.Sum(nil)[:]
	return bs
}
func Sha512(input []byte) []byte {
	hmac512 := hmac.New(sha512.New, []byte("TossP.com"))
	hmac512.Write(input)
	bs := hmac512.Sum(nil)[:]
	return bs
}

func Hash32(password, salt []byte) []byte {
	return pbkdf2.Key(password, Sha256(salt), 100, 32, sha256.New)[:]
}

func HashKey(input []byte, keylen int) (key []byte) {
	return pbkdf2.Key(input, Sha1([]byte("TossP.com")), 1024, keylen, sha256.New)[:]
}

func GmHashKey(input []byte, keylen int) (key []byte) {
	//return sm3.Sm3Sum(input)[:]
	return pbkdf2.Key(sm3.Sm3Sum(input), Sha1([]byte("TossP.com")), 1024, keylen, sha256.New)[:]
}
func HashSha(input, salt []byte, keylen int) (key []byte) {
	return pbkdf2.Key(input, Sha512(salt), 1024, keylen, sha512.New)[:]
}
func HashArgon(input, salt []byte, keylen uint32) (key []byte) {
	return argon2.IDKey(input, Sha512(salt), 4, 64*1024, 4, keylen)
}

func HashFile(filename string) (bs []byte, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer func() {
		_ = file.Close()
	}()
	bs, err = HashReader(file)
	return
}

func HashReader(file io.Reader) (bs []byte, err error) {
	hash := hmac.New(sha512.New, []byte("TossP.com"))
	if _, err = io.Copy(hash, file); err != nil {
		return
	}
	bs = hash.Sum(nil)[:]
	return
}

func inStrArray(val string, array []string) bool {
	for _, item := range array {
		if item == val {
			return true
		}
	}
	return false
}
