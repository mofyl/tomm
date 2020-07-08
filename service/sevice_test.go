package service

import (
	"testing"
	"time"
	"tomm/api/service"
	"tomm/log"
	"tomm/utils"
)

func getAESBaseStr() string {

	secretKey := "2ffd7fbe21a5e6eb3321d723900a79f0"
	//appKey := "055285a69ec81f6477e49fe95da22eba"
	sendTime := time.Now().Unix()
	dataInfo := service.ReqDataInfo{
		SendTime:   sendTime,
		ExtendInfo: nil,
	}

	data, _ := utils.Json.Marshal(dataInfo)

	base64Str, err := utils.AESCBCBase64Encode(secretKey, data)
	if err != nil {
		log.Error("AESCBCBase64Encode Fail err is %s", err.Error())
		return ""
	}

	log.Info("Req Base64Str is %s", base64Str)
	origData, err := utils.AESCBCBase64Decode(secretKey, base64Str)

	if err != nil {
		log.Info("Decode err is %s", err)
		return ""
	}
	res := service.ReqDataInfo{}
	err = utils.Json.Unmarshal(origData, &res)
	if err != nil {
		log.Info("Json Unmarshal err is %s", err)
		return ""
	}

	if res.SendTime != sendTime {
		panic("AESBase Fail")
	}

	log.Info("Get Data is %s", res.SendTime)
	return base64Str
}

func TestGetToken(t *testing.T) {
	str := getAESBaseStr()
	if str == "" {
		panic("Get AESBaseStr Fail")
	}

}

//
//func TestGetToken(t *testing.T) {
//
//	secretKey := "2ffd7fbe21a5e6eb3321d723900a79f0"
//
//	getTokenRes := service.GetTokenRes{}
//
//	getTokenRes.Token = "9IzGNxXoPdL3hoYhnlj5Ag0LBvvgNnd4n10o1LgZJAUQxY8aQQ6CIXG5pqIgFOnzUdMPqtR4mSK9VK5PnqNicA"
//
//	tokenInfoB, err := utils.AESCBCBase64Decode(secretKey, getTokenRes.Token)
//	if err != nil {
//		log.Info("Decode Fail Err is %s", err.Error())
//		return
//	}
//	info := service.TokenInfo{}
//
//	utils.Json.Unmarshal(tokenInfoB, &info)
//
//	log.Info("Token is %s , ExpTime is %d", info.Token, info.ExpTime)
//}
