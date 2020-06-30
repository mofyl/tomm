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

func TestDecode(t *testing.T) {

	str := "laX2AhoKH1mTm1XO/dWD8tIABgEW0m2sYVfht9y5CRkNSAdncYspMI54XCixSwi7Ef+qEh0dp7KcGKjE5AHlXQ== "
	key := "2ffd7fbe21a5e6eb3321d723900a79f0"
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
