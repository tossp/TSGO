package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"io"
	"math/big"
)

//NewKey 随机生成密钥
func NewKey() (prk *ecdsa.PrivateKey, err error) {
	return GenerateKey(rand.Reader)
}

//NewKeyWithKey 指定特征SHA512后生成密钥
func NewKeyWithKey(key []byte) (prk *ecdsa.PrivateKey) {
	hash := HashKey(key, P521().Params().BitSize)
	prk, _ = MakeKeyWithKey(hash)
	return
}

//MakeKeyWithKey 指定特征后直接生成密钥
func MakeKeyWithKey(key []byte) (prk *ecdsa.PrivateKey, err error) {
	return GenerateKey(bytes.NewReader(key))
}

//GenerateKey 生成密钥对
func GenerateKey(r io.Reader) (prk *ecdsa.PrivateKey, err error) {
	//prk, err := ecdsa.GenerateKey(curve, strings.NewReader(randKey))

	if prk, err = ecdsa.GenerateKey(P521(), r); err != nil {
		fmt.Printf("Crypt 初始化失败 %s\n 需要 = %d", err, P521().Params().BitSize)
	}
	return
}

// 私钥 -> []byte
func FromECDSA(priv *ecdsa.PrivateKey) []byte {
	if priv == nil {
		return nil
	}
	return PaddedBigBytes(priv.D, priv.Params().BitSize/8+8)
	//return priv.D.Bytes()
}

// []byte -> 私钥
func ToECDSA(d []byte) (*ecdsa.PrivateKey, error) {
	return toECDSA(d, false)
}

func toECDSA(d []byte, strict bool) (*ecdsa.PrivateKey, error) {
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = P521()
	lenD := 8*(len(d)-8) + 1
	if strict && lenD != priv.Params().BitSize {
		return nil, fmt.Errorf("invalid length, need %d bits, %d %d", priv.Params().BitSize, len(d), lenD)
	}
	priv.D = new(big.Int).SetBytes(d)

	// The priv.D must < N
	if priv.D.Cmp(P521().Params().N) >= 0 {
		return nil, fmt.Errorf("invalid private key, >=N")
	}
	// The priv.D must not be zero or negative.
	if priv.D.Sign() <= 0 {
		return nil, fmt.Errorf("invalid private key, zero or negative")
	}

	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(d)
	if priv.PublicKey.X == nil {
		return nil, fmt.Errorf("invalid private key")
	}
	return priv, nil
}

// 公钥 -> []byte
func FromECDSAPub(pub *ecdsa.PublicKey) []byte {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}
	return elliptic.Marshal(P521(), pub.X, pub.Y)
}

// []byte -> 公钥
func ToECDSAPub(pub []byte) *ecdsa.PublicKey {
	if len(pub) == 0 {
		return nil
	}
	x, y := elliptic.Unmarshal(P521(), pub)
	return &ecdsa.PublicKey{Curve: P521(), X: x, Y: y}
}

// GenerateSharedSecret 生成共享密钥
func GenerateSharedSecret(privKey *ecdsa.PrivateKey, pubKey *ecdsa.PublicKey) []byte {
	x, _ := P521().ScalarMult(pubKey.X, pubKey.Y, privKey.D.Bytes())
	k := Sha512(x.Bytes())[:]
	return k[:]
}

func M(privKey *ecdsa.PrivateKey) {
	x509.MarshalECPrivateKey(privKey)
}

func P521() elliptic.Curve {
	return elliptic.P521()
}
