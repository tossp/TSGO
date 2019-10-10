package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math/big"
)

//Sign 签名
func Sign(priv *ecdsa.PrivateKey, pt []byte) (sign []byte, err error) {
	// 根据明文plaintext和私钥，生成两个big.Ing
	r, s, err := ecdsa.Sign(rand.Reader, priv, pt)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	rs, err := r.MarshalText()
	if err != nil {
		return nil, err
	}
	ss, err := s.MarshalText()
	if err != nil {
		return nil, err
	}
	// 将r，s合并（以“+”分割），作为签名返回
	var b bytes.Buffer
	b.Write(rs)
	b.Write([]byte(`+`))
	b.Write(ss)
	return b.Bytes(), nil
}

//SignVer 签名验证
func SignVer(pub *ecdsa.PublicKey, pt, sign []byte) (ok bool) {
	var rint, sint big.Int
	rs := bytes.Split(sign, []byte("+"))
	if err := rint.UnmarshalText(rs[0]); err != nil {
		return
	}
	if err := sint.UnmarshalText(rs[1]); err != nil {
		return
	}
	// 根据公钥，明文，r，s验证签名
	ok = ecdsa.Verify(pub, pt, &rint, &sint)
	return
}
