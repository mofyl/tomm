package service

import (
	"encoding/json"
	"time"
	"tomm/core/server"
	"tomm/ecode"
	"tomm/log"
	"tomm/service/oauth"
	"tomm/utils"
)

type Ser struct {
	e *server.Engine
}

func NewService() *Ser {
	s := &Ser{}

	e := server.NewEngine(nil)

	s.e = e
	s.registerRouter()
	return s
}

func (s *Ser) registerRouter() {
	s.e.GET("/getToken", s.getToken)
	s.e.GET("/verifyToken", s.verifyToken)
}

func (s *Ser) Close() {
	s.e.Close()
}

func (s *Ser) Start() {
	s.e.RunServer()
}

func (s *Ser) verifyToken(c *server.Context) {
	req := VerifyTokenReq{}
	err := c.Bind(&req)
	if err != nil {
		log.Warn("VerifyToken Bind Err is %s", err.Error())
	}

	// 查看token是否存在
	token, expTime, err := oauth.GetToken(req.AppKey)

	if err != nil {
		log.Error("Verify Token Fail Err is %s", err.Error())
		httpCode(c, ecode.ParamFail)
	}

	if token != token {
		// 校验失败
		httpCode(c, ecode.VerifyFail)
	}
	res := VerifyTokenRes{}
	res.ExpTime = expTime
	httpData(c, res)
}

func (s *Ser) getToken(c *server.Context) {
	req := GetTokenReq{}
	err := c.Bind(&req)

	if err != nil {
		log.Warn("GetToken Bind Err is %s ", err.Error())

		httpCode(c, ecode.ParamFail)
		return
	}

	if req.AppKey == "" || req.Data == "" {
		httpCode(c, ecode.ParamFail)
		return
	}

	// 获取该appKey
	secretInfo, err := oauth.GetOAuthInfo(req.AppKey)
	if err != nil || secretInfo == nil {
		if err != nil {
			log.Error("GetToken Fail AppKey is %s , Err is %s", req.AppKey, err.Error())
		}
		httpCode(c, ecode.SecretKeyFail)
		return
	}

	reqDataInfo, eCode := GetDataInfo(secretInfo.SecretKey, req.Data)
	if eCode != nil {
		log.Error("Get Data Info Fail ")
		httpCode(c, eCode)
		return
	}
	// 超过10分钟就不处理了
	if time.Now().Unix()-int64(reqDataInfo.SendTime) > MAX_TTL {
		httpCode(c, ecode.PackageTimeout)
		log.Error("Package Timeout ")
		return
	}

	token, expTime, err := oauth.GetToken(req.AppKey)
	if err != nil {
		httpCode(c, ecode.SystemErr)
		log.Error("Get Token Fail ")
		return
	}
	// 更新ChannelInfo
	if token != "" && reqDataInfo.ChannelInfo != "" {
		secretInfo.ChannelInfo = reqDataInfo.ChannelInfo
		oauth.UpdateChannelInfo(secretInfo)
	}

	log.Debug("Return Token is %s", token)
	// 构造返回值
	// 返回值包括 token + expTime + extendInfo
	tokenInfo := TokenInfo{
		Token:      token,
		ExpTime:    expTime,
		ExtendInfo: reqDataInfo.ExtendInfo,
	}
	tokenB, err := json.Marshal(tokenInfo)
	resBase64Str, err := utils.AESCBCBase64Encode(secretInfo.SecretKey, tokenB)
	if err != nil {
		log.Error("AESCBCBase64Encode Fail Err is %s", err.Error())
		httpCode(c, ecode.SystemErr)
		return
	}
	res := GetTokenRes{}
	res.TokenInfo = resBase64Str
	httpData(c, res)
	log.Error("Cur Res is %v", res)
}

//
//func buildRes(c *server.Context, res *GetTokenRes, secretKey string, code int) {
//	var b []byte
//	if res != nil {
//		b, _ := json.Marshal(res)
//		origData, err := utils.AESCBCBase64Encode(secretKey, b)
//		if err != nil {
//			code = 500
//			b = nil
//		} else {
//			b = []byte(origData)
//		}
//	} else {
//		code = 400
//		b = nil
//	}
//
//	c.Json(code, b)
//}

func GetDataInfo(secretKey string, data string) (*ReqDataInfo, ecode.ErrMsgs) {

	// 使用 secretKey 进行 AES解密
	origData, err := utils.AESCBCBase64Decode(secretKey, data)

	if err != nil {
		return nil, ecode.DeCodeFail
	}

	origLen := len(origData)

	if origLen < DATALEN+TIMELEN+CHANNEL_INFO_LEN || origLen > MAX_DATA {
		return nil, ecode.ParamFail
	}
	//
	//dataLen := binary.BigEndian.Uint32(origData[:DATALEN])
	//sendTime := binary.BigEndian.Uint64(origData[DATALEN : DATALEN+TIMELEN])
	//channelInfo := origData[DATALEN+TIMELEN : DATALEN+TIMELEN+CHANNEL_INFO_LEN]
	//extendInfo := origData[DATALEN+TIMELEN+CHANNEL_INFO_LEN:]
	//reqInfo := &ReqDataInfo{
	//	DataLen:     int32(dataLen),
	//	SendTime:    int64(sendTime),
	//	ChannelInfo: string(channelInfo),
	//	ExtendInfo:  extendInfo,
	//}

	reqInfo := &ReqDataInfo{}
	err = json.Unmarshal(origData, reqInfo)

	if err != nil {
		return nil, ecode.ParamFail
	}
	//if int(reqInfo.DataLen) != origLen {
	//	return nil, DECODE_FAIL, DECODE_FAIL_MSG
	//}
	return reqInfo, nil
}
