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

func TestCheckUrl(t *testing.T) {

	//url := "https://www.baiduqzxc232.com.cn"
	url := "1.0.0.1.1"

	if CheckUrl(url) {
		fmt.Println("Success")
	} else {
		fmt.Println("Fail")
	}
}

func Test_test(t *testing.T) {
	var m = map[string]int{
		"A": 21,
		"B": 22,
		"C": 23,
	}

	counter := 0
	for k, _ := range m {
		if counter == 0 {
			delete(m, "A")
		}
		counter++
		fmt.Println(k)
	}
	fmt.Println(counter)
}
