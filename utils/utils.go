package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"strings"
)

const (
	PRO_NAME = "tomm"
)

func GetProDirAbs() string {
	sbuilder := strings.Builder{}
	sbuilder.WriteString(os.Getenv("GOPATH"))
	sbuilder.WriteString(string(filepath.Separator))
	sbuilder.WriteString("src")
	sbuilder.WriteString(string(filepath.Separator))
	sbuilder.WriteString(PRO_NAME)
	sbuilder.WriteString(string(filepath.Separator))

	path := sbuilder.String()
	return path
}

func GetUUID() (uuid.UUID, error) {
	return uuid.NewUUID()
}

func AESCBCEncode(key string, data []byte) ([]byte, error) {
	keyb := []byte(key)
	block, err := aes.NewCipher(keyb)
	if err != nil {
		return nil, err
	}
	content := pKCS5Padding(data, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, keyb[:block.BlockSize()])
	crypted := make([]byte, len(content))
	blockMode.CryptBlocks(crypted, content)
	return crypted, nil
}

func AESCBCDecode(key string, data []byte) ([]byte, error) {
	keyb := []byte(key)
	block, err := aes.NewCipher(keyb)

	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, keyb[:blockSize])
	origData := make([]byte, len(data))
	blockMode.CryptBlocks(origData, data)

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
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func Base64Encode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func Base64Decode(data string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(data)
}

func AESCBCBase64Encode(key string, data []byte) (string, error) {
	crypted, err := AESCBCEncode(key, data)

	if err != nil {
		return "", err
	}

	return Base64Encode(crypted), nil
}

func AESCBCBase64Decode(key string, data string) ([]byte, error) {
	crypted, err := Base64Decode(data)

	if err != nil {
		return nil, err
	}

	orig, err := AESCBCDecode(key, crypted)
	if err != nil {
		return nil, err
	}
	return orig, nil
}

func StrUUID() (string, error) {

	uuid, err := GetUUID()
	if err != nil {
		return "", nil
	}

	uuidB, _ := uuid.MarshalText()
	res := md5.Sum(uuidB)
	return hex.EncodeToString(res[:]), nil
}
