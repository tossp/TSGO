package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
)

func EccEncrypt(priv *ecdsa.PrivateKey, pub *ecdsa.PublicKey, plantText []byte) []byte {
	return AesEncrypt(plantText, GenerateSharedSecret(priv, pub))
}

func EccDecrypt(priv *ecdsa.PrivateKey, pub *ecdsa.PublicKey, cipherText []byte) []byte {
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

func AesDecrypt(cipherText, key []byte) []byte {
	k := HashKey(key, 32)
	block, _ := aes.NewCipher(k) //选择加密算法
	blockModel := cipher.NewCBCDecrypter(block, k[:block.BlockSize()])

	blockModel.CryptBlocks(cipherText, cipherText)
	plantText := PKCS7UnPadding(cipherText)
	return plantText
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
