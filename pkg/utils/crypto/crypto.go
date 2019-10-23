package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"errors"
)

func EccEncrypt(priv *ecdsa.PrivateKey, pub *ecdsa.PublicKey, plantText []byte) []byte {
	return AesEncrypt(plantText, GenerateSharedSecret(priv, pub))
}

func EccDecrypt(priv *ecdsa.PrivateKey, pub *ecdsa.PublicKey, cipherText []byte) ([]byte, error) {
	return AesDecrypt(cipherText, GenerateSharedSecret(priv, pub))
}

func AesEncrypt(plantText, key []byte) []byte {
	k := HashKey(key, 32)
	block, _ := aes.NewCipher(k) //选择加密算法
	plantText = PKCS7Padding(plantText, block.BlockSize())
	blockModel := cipher.NewCBCEncrypter(block, k[:block.BlockSize()])
	ciphertext := make([]byte, len(plantText))
	blockModel.CryptBlocks(ciphertext, plantText)
	return ciphertext
}

func AesDecrypt(cipherText, key []byte) (plantText []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			//check exactly what the panic was and create error.
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknow panic")
			}
		}

	}()

	k := HashKey(key, 32)
	block, _ := aes.NewCipher(k) //选择加密算法
	blockModel := cipher.NewCBCDecrypter(block, k[:block.BlockSize()])

	blockModel.CryptBlocks(cipherText, cipherText)
	plantText = PKCS7UnPadding(cipherText)
	return
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}
