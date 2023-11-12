package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

const (
	key = "cvrhsdftredhghgfjhgwsfresdsfhjk"
)

func getKeyBytes(key string) []byte {
	keyBytes := []byte(key)
	switch l := len(keyBytes); {
	case l < 16:
		keyBytes = append(keyBytes, make([]byte, 16-l)...)
	case l > 16:
		keyBytes = keyBytes[:16]
	}
	return keyBytes
}

func encrypt(key string, origData []byte) ([]byte, error) {
	keyBytes := getKeyBytes(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = pKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, keyBytes[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func decrpt(key string, crypted []byte) ([]byte, error) {
	keyBytes := getKeyBytes(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, keyBytes[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = pKCS5UnPadding(origData)
	return origData, nil
}

func pKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func Encrypt(val string) (string, error) {
	m := md5.Sum([]byte(key))
	aesKey := hex.EncodeToString(m[:])[0:16]
	origData := []byte(val)
	crypted, err := encrypt(aesKey, origData)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(crypted), nil
}

func Decrypt(val string) (string, error) {
	m := md5.Sum([]byte(key))
	aesKey := hex.EncodeToString(m[:])[0:16]
	base64Data, _ := base64.StdEncoding.DecodeString(val)
	origData, err := decrpt(aesKey, base64Data)
	if err != nil {
		return "", err
	}
	return string(origData), nil
}
