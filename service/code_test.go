package service

import (
	"encoding/binary"
	"fmt"
	"testing"
	"time"
	"tomm/log"
	"tomm/utils"
)

var (
	key = "2ffd7fbe21a5e6eb3321d723900a79f0"
)

func TestCode(t *testing.T) {
	str, err := buildData()

	if err != nil {
		log.Error("TestCode BuildData Fail %s", err.Error())
		return
	}
	log.Info(" BuildData Str %s", str)

	//data, errCode := getDataInfo(key, str)
	//
	//if errCode != nil {
	//	return
	//}
	//
	//log.Debug("TestCode Data Code is %d", data.TimeStamp)
	//log.Debug("TestCode Data Code is %s", data.Code)

}

func buildData() (string, error) {

	code := "a3f45aeedc77954784c4a62adb0a3255"

	buf := make([]byte, CODE_TIME_LEN, CODE_DATA_LEN)

	binary.BigEndian.PutUint64(buf, uint64(time.Now().Unix()))

	buf = append(buf, []byte(code)...)

	return utils.AESCBCBase64Encode(key, buf)
	//
	//select {
	//case channel <- 1:
	//default:
	//
	//}
}

func Test_test(t *testing.T) {

	channel := make(chan int, 1)

	select {
	case channel <- 1:
	default:
		fmt.Println(111)
	}
}
