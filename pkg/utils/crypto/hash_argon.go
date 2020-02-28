package crypto

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"runtime"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argon2d = iota
	argon2i
	argon2id
)
const (
	argonDefaultMemoryPasses = 4
	argonDefaultMemorySize   = 64 * 1024
	argonDefaultHashSize     = 32
	argonDefaultSaltLength   = 16
	argonDefaultMode         = argon2id
	argonParametersFormat    = "m=%d,i=%d,s=%d,p=%d,k=%d,l=%d"
)

var (
	argonModi = []string{"i", "id"} // argon modus

	ErrBadParameters       = errors.New("哈希参数格式错误")
	ErrIncompatibleVersion = errors.New("版本不兼容")
	ErrInvalidHash         = errors.New("编码哈希的格式不正确")
)

type argon2Param struct {
	iterations  uint32 // time setting
	memorySize  uint32 // memory setting in KiB, e.g. 64*1024 -> 64MB
	parallelism uint8  // threads setting
	mode        uint8  // modus for argon, i or id
	keyLength   uint32 // hash size in bytes (min. 16)
	saltLength  uint32 // salt size in bytes
}

func init() {
	h := &argon2Param{
		iterations:  argonDefaultMemoryPasses,
		memorySize:  argonDefaultMemorySize,
		parallelism: uint8(runtime.NumCPU() / 2),
		mode:        argonDefaultMode,
		keyLength:   argonDefaultHashSize,
		saltLength:  argonDefaultSaltLength,
	}
	implementations[h.GetID()] = h
}

func (h *argon2Param) GetID() string {
	return "argon2"
}
func (h *argon2Param) String() string {
	return h.encodeConfigure()
}

func (h *argon2Param) Generate(plainText string, param ...string) (encodedHash string, err error) {
	salt, err := GenerateRandomBytes(int(h.saltLength))
	if err != nil {
		return
	}
	hash := h.hash([]byte(plainText), salt)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	encodedHash = fmt.Sprintf("$%s$v=%d$%s$%s$%s", h.GetID(), argon2.Version, h.encodeConfigure(), b64Salt, b64Hash)
	return encodedHash, nil
}

func (h *argon2Param) Compare(plainText, encodedHash string) (match bool, err error) {
	pp, salt, hash, err := h.decodeHash(encodedHash)
	if err != nil {
		return false, err
	}
	otherHash := pp.hash([]byte(plainText), salt)
	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func (h *argon2Param) NeedsReHash(encodedHash string) (bool, error) {
	pp, _, _, err := h.decodeHash(encodedHash)
	if err != nil {
		return false, err
	}
	return pp.encodeConfigure() != h.encodeConfigure(), nil
}

func (h *argon2Param) Configure(param string) (hash, error) {
	p := new(argon2Param)
	if _, err := p.decodeConfigure(param); err != nil {
		return nil, err
	}

	if p.mode > argon2id || p.keyLength < 16 || p.iterations <= 0 || p.memorySize <= 0 {
		return nil, ErrBadParameters
	}
	return p, nil
}

func (h *argon2Param) hash(plain, salt []byte) (hash []byte) {
	switch h.mode {
	case argon2i:
		hash = argon2.Key(plain, salt, h.iterations, h.memorySize, h.parallelism, h.keyLength)
	default:
		hash = argon2.IDKey(plain, salt, h.iterations, h.memorySize, h.parallelism, h.keyLength)
	}
	return
}
func (h *argon2Param) decodeHash(encodedHash string) (param *argon2Param, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}
	if h.GetID() != vals[1] {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	param = &argon2Param{}
	_, err = param.decodeConfigure(vals[3])
	if err != nil {
		return nil, nil, nil, ErrBadParameters
	}
	salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, ErrBadParameters
	}
	if param.saltLength != uint32(len(salt)) {
		return nil, nil, nil, ErrBadParameters
	}
	hash, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, ErrBadParameters
	}
	if param.keyLength != uint32(len(hash)) {
		return nil, nil, nil, ErrBadParameters
	}
	return
}

func (h *argon2Param) decodeConfigure(param string) (n int, err error) {
	n, err = fmt.Sscanf(param, argonParametersFormat, &h.mode, &h.iterations, &h.memorySize, &h.parallelism, &h.keyLength, &h.saltLength)
	if err != nil {
		return
	}
	return
}
func (h *argon2Param) encodeConfigure() string {
	return fmt.Sprintf(argonParametersFormat, h.mode, h.iterations, h.memorySize, h.parallelism, h.keyLength, h.saltLength)
}
