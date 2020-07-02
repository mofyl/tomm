package service

import (
	"testing"
	"time"
	"tomm/log"
	"tomm/utils"
)

func TestEncodeDecode(t *testing.T) {

	secretKey := "2ffd7fbe21a5e6eb3321d723900a79f0"
	//appKey := "055285a69ec81f6477e49fe95da22eba"

	dataInfo := ReqDataInfo{
		SendTime:    time.Now().Unix(),
		ChannelInfo: "abc",
		ExtendInfo:  nil,
	}

	data, _ := utils.Json.Marshal(dataInfo)

	base64Str, err := utils.AESCBCBase64Encode(secretKey, data)
	if err != nil {
		log.Error("AESCBCBase64Encode Fail err is %s", err.Error())
		return
	}

	log.Info("Req Base64Str is %s", base64Str)
	origData, err := utils.AESCBCBase64Decode(secretKey, base64Str)

	if err != nil {
		log.Info("Decode err is %s", err)
		return
	}
	res := ReqDataInfo{}
	err = utils.Json.Unmarshal(origData, &res)
	if err != nil {
		log.Info("Json Unmarshal err is %s", err)
		return
	}

	log.Info("Get Data is %s", res.ChannelInfo)
}

func TestGetToken(t *testing.T) {

	secretKey := "2ffd7fbe21a5e6eb3321d723900a79f0"
	//appKey := "055285a69ec81f6477e49fe95da22eba"

	//dataInfo := ReqDataInfo{
	//	SendTime:    time.Now().Unix(),
	//	ChannelInfo: "abc",
	//	ExtendInfo:  nil,
	//}
	//
	//data, _ := json.Marshal(dataInfo)
	//
	//base64Str, err := utils.AESCBCBase64Encode(secretKey, data)
	//if err != nil {
	//	log.Error("AESCBCBase64Encode Fail err is %s", err.Error())
	//	return
	//}
	//log.Info("Req Base64Str is %s ", base64Str)
	//// 开始客户端请求
	//resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:8086/getToken?app_key=%s&data=%s", appKey, base64Str))
	//
	//if err != nil {
	//	log.Error("http Get Fail err is %s", err.Error())
	//	return
	//}
	//defer resp.Body.Close()
	////buf := bytes.NewBuffer([]byte{})
	//body, err := ioutil.ReadAll(resp.Body)
	////json.NewDecoder(buf).
	////n, err := resp.Body.Read(body)
	//if err != nil {
	//	log.Error("Resp read body Fail err is %s", err.Error())
	//	return
	//}
	getTokenRes := GetTokenRes{}
	//err = json.Unmarshal(body, &getTokenRes)
	//if err != nil {
	//	log.Error("unmarshal body Fail err is %s", err.Error())
	//	return
	//}
	getTokenRes.TokenInfo = "9IzGNxXoPdL3hoYhnlj5Ag0LBvvgNnd4n10o1LgZJAUQxY8aQQ6CIXG5pqIgFOnzUdMPqtR4mSK9VK5PnqNicA"
	//log.Info("GetToken Res is ErrCode %d , ErrMsg %s , TokenInfo is %s", getTokenRes.ErrCode, getTokenRes.ErrMsg, getTokenRes.TokenInfo)

	tokenInfoB, err := utils.AESCBCBase64Decode(secretKey, getTokenRes.TokenInfo)
	if err != nil {
		log.Info("Decode Fail Err is %s", err.Error())
		return
	}
	info := TokenInfo{}

	utils.Json.Unmarshal(tokenInfoB, &info)

	log.Info("Token is %s , ExpTime is %d", info.Token, info.ExpTime)
}
