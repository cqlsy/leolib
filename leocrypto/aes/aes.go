package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"errors"
)

func AesEncrypt(key []byte, data []byte) ([]byte, error) {
	keyHash := sha512.Sum512(key)
	aesKey := keyHash[:32]
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	cbcIv := make([]byte, 16)
	rand.Read(cbcIv)
	data = PKCS5Padding(data, block.BlockSize())
	crypted := make([]byte, len(data)+16)
	copy(crypted[:16], cbcIv)
	blockMode := cipher.NewCBCEncrypter(block, cbcIv)
	blockMode.CryptBlocks(crypted[16:], data)
	return []byte(base64.StdEncoding.EncodeToString(crypted)), nil
}

func AesDecrypt(key []byte, data []byte) (origData []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New("AES Decrypt error")
		}
	}()
	data, _ = base64.StdEncoding.DecodeString(string(data))
	cbcIv := data[:16]
	keyHash := sha512.Sum512(key)
	aesKey := keyHash[:32]
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, cbcIv)
	origData = make([]byte, len(data[16:]))
	blockMode.CryptBlocks(origData, data[16:])
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
