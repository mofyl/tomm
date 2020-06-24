package utils

import (
	"fmt"
	"testing"
)

func TestAESCBCBase64(t *testing.T) {
	key := "qwertyuiopasdfgh"
	data := []byte("hello world")

	str, err := AESCBCBase64Encode(key, data)

	if err != nil {
		fmt.Println("Encode err ", err.Error())
		return
	}

	fmt.Println("Encode Res ", str)
	orig, err := AESCBCBase64Decode(key, str)
	if err != nil {
		fmt.Println("Decode err ", err.Error())
		return
	}

	fmt.Println(string(orig))
}

func TestAESCBC(t *testing.T) {
	key := "qwertyuiopasdfgh"
	data := []byte("hello world")

	crypted, err := AESCBCEncode(key, data)
	if err != nil {
		fmt.Println("Encode err ", err.Error())
		return
	}

	origData, err := AESCBCDecode(key, crypted)
	if err != nil {
		fmt.Println("Decode err ", err.Error())
		return
	}

	fmt.Println(string(origData))
}
