package crypto

import (
	"bytes"
	"crypto/elliptic"
	"crypto/rand"
	"io"
	"math/big"

	"github.com/tjfoc/gmsm/sm2"
)

var one = new(big.Int).SetInt64(1)

//NewKey 随机生成密钥
func NewSm2Key() (prk *sm2.PrivateKey, err error) {
	return Sm2GenerateKey(rand.Reader)
}

//NewKeyWithKey 指定特征SHA512后生成密钥
func NewSm2KeyWithKey(key []byte) (prk *sm2.PrivateKey) {
	hash := HashKey(key, P256Sm2().Params().BitSize/8+8)
	prk, _ = MakeSm2KeyWithKey(hash)
	return
}

//MakeKeyWithKey 指定特征后直接生成密钥
func MakeSm2KeyWithKey(key []byte) (prk *sm2.PrivateKey, err error) {
	return Sm2GenerateKey(bytes.NewReader(key))
}

func Sm2GenerateKey(r io.Reader) (*sm2.PrivateKey, error) {
	c := P256Sm2()
	k, err := randFieldElement(c, r)
	if err != nil {
		return nil, err
	}
	priv := new(sm2.PrivateKey)
	priv.PublicKey.Curve = c
	priv.D = k
	priv.PublicKey.X, priv.PublicKey.Y = c.ScalarBaseMult(k.Bytes())
	return priv, nil
}

func randFieldElement(c elliptic.Curve, rand io.Reader) (k *big.Int, err error) {
	params := c.Params()
	b := make([]byte, params.BitSize/8+8)
	_, err = io.ReadFull(rand, b)
	if err != nil {
		return
	}
	k = new(big.Int).SetBytes(b)
	n := new(big.Int).Sub(params.N, one)
	k.Mod(k, n)
	k.Add(k, one)
	return
}

func FromsSm2Pub(pub *sm2.PublicKey) []byte {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}
	return elliptic.Marshal(P256Sm2(), pub.X, pub.Y)
}

// []byte -> 公钥
func ToSm2Pub(pub []byte) *sm2.PublicKey {
	if len(pub) == 0 {
		return nil
	}
	x, y := elliptic.Unmarshal(P256Sm2(), pub)
	return &sm2.PublicKey{Curve: P256Sm2(), X: x, Y: y}
}

// GenerateSm2SharedSecret 生成Sm2共享密钥
func GenerateSm2SharedSecret(privKey *sm2.PrivateKey, pubKey *sm2.PublicKey) (key []byte) {
	x, _ := P256Sm2().ScalarMult(pubKey.X, pubKey.Y, privKey.D.Bytes())
	key = x.Bytes()
	return
}

func P256Sm2() elliptic.Curve {
	return sm2.P256Sm2()
}
