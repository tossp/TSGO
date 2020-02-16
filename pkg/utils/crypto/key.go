package crypto

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"io"
	"math/big"

	"github.com/tjfoc/gmsm/sm2"
	"github.com/tossp/tsgo/pkg/errors"
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
func GenerateSharedSecret(priv crypto.PrivateKey, pub crypto.PublicKey, salt ...byte) ([]byte, error) {
	var (
		x1    *big.Int
		y1    *big.Int
		k     []byte
		curve elliptic.Curve
	)
	switch key := priv.(type) {
	case *ecdsa.PrivateKey:
		k = key.D.Bytes()
		pubKey, ok := pub.(*ecdsa.PublicKey)
		if !ok {
			fmt.Println("pub only support ecdsa.PublicKey point type")
			return nil, errors.New("pub only support ecdsa.PublicKey point type")
		}
		x1 = pubKey.X
		y1 = pubKey.Y
		curve = pubKey.Curve
	case *sm2.PrivateKey:
		k = key.D.Bytes()
		pubKey, ok := pub.(*sm2.PublicKey)
		if !ok {
			fmt.Println("pub only support sm2.PublicKey point type")
			return nil, errors.New("pub only support sm2.PublicKey point type")
		}
		x1 = pubKey.X
		y1 = pubKey.Y
		curve = pubKey.Curve
	default:
		fmt.Println("priv only support ecdsa.PrivateKey and sm2.PrivateKey")
		return nil, errors.New("priv only support ecdsa.PrivateKey and sm2.PrivateKey")
	}

	if salt == nil {
		salt = []byte("TossP.com")
	}
	x, _ := curve.ScalarMult(x1, y1, k)
	key := Sha512(append(x.Bytes(), salt...))[:]
	return key[:], nil
}

func M(privKey *ecdsa.PrivateKey) {
	x509.MarshalECPrivateKey(privKey)
}

func P521() elliptic.Curve {
	return elliptic.P521()
}
