package service

import (
	"encoding/json"
	"go.uber.org/zap"
	"time"
	"tomm/core/server"
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
	res := VerifyTokenRes{}
	if err != nil {
		log.Warn("VerifyToken", zap.String("Bind Error ", err.Error()))
	}

	// 查看token是否存在
	token, expTime, err := oauth.GetToken(req.AppKey)

	if err != nil {
		log.Error("Verify Token Fail", zap.String("error msg is ", err.Error()))
		res.ErrCode = VERIFY_FAIL
		res.ErrMsg = VERIFY_FAIL_MSG
		c.Json(200, res)
	}

	if token != token {
		// 校验失败
		res.ErrMsg = VERIFY_FAIL_MSG
		res.ErrCode = VERIFY_FAIL
		c.Json(200, res)
	}
	res.ErrCode = SUCCESS
	res.ErrMsg = SUCCESS_MSG
	res.ExpTime = expTime
	c.Json(200, res)
}

func (s *Ser) getToken(c *server.Context) {
	req := GetTokenReq{}
	err := c.Bind(&req)
	res := GetTokenRes{}
	if err != nil {
		log.Warn("GetToken", zap.String("Bind Error ", err.Error()))
		res.ErrCode = PARAM_FAIL
		res.ErrMsg = PARAM_FAIL_MSG
		c.Json(200, res)
		return
	}

	if req.AppKey == "" || req.Data == "" {
		res.ErrMsg = PARAM_FAIL_MSG
		res.ErrCode = PARAM_FAIL
		c.Json(200, res)
		return
	}

	// 获取该appKey
	secretInfo, err := oauth.GetOAuthInfo(req.AppKey)
	if err != nil || secretInfo == nil {
		if err != nil {
			log.Warn("GetToken", zap.String("Redis Get Key Fail AppKey is ", req.AppKey),
				zap.String("Redis Key Get Fail err is ", err.Error()))
		}
		res.ErrCode = SECRET_KEY_FAIL
		res.ErrMsg = SECRET_KEY_FAIL_MSG
		c.Json(200, res)
		return
	}
	// 更新ChannelInfo
	if secretInfo.ChannelInfo != "" {
		oauth.UpdateChannelInfo(secretInfo)
	}

	reqDataInfo, errCode, errMsg := GetDataInfo(secretInfo.SecretKey, req.Data)
	if errCode != SUCCESS {
		res.ErrCode = errCode
		res.ErrMsg = errMsg
		c.Json(200, res)
	}
	// 超过10分钟就不处理了
	if time.Now().Unix()-int64(reqDataInfo.SendTime) > MAX_TTL {
		res.ErrCode = PACKAGE_TIME_OUT
		res.ErrMsg = PACKAGE_TIME_OUT_MSG
		c.Json(200, res)
		return
	}

	token, expTime, err := oauth.GetToken(req.AppKey)
	if err != nil {
		res.ErrMsg = SYSTEM_FAILE_MSG
		res.ErrCode = SYSTEM_FAIL
		c.Json(200, res)
		return
	}

	// 构造返回值
	// 返回值包括 token + expTime + extendInfo
	res.Token = token
	res.ExpTime = expTime
	res.ExtendInfo = reqDataInfo.ExtendInfo
	res.ErrMsg = SUCCESS_MSG
	res.ErrCode = SUCCESS
	c.Json(200, res)
}

func buildRes(c *server.Context, res *GetTokenRes, secretKey string, code int) {
	var b []byte
	if res != nil {
		b, _ := json.Marshal(res)
		origData, err := utils.AESCBCBase64Encode(secretKey, b)
		if err != nil {
			code = 500
			b = nil
		} else {
			b = []byte(origData)
		}
	} else {
		code = 400
		b = nil
	}

	c.Json(code, b)
}

func GetDataInfo(secretKey string, data string) (*ReqDataInfo, int64, string) {

	// 使用 secretKey 进行 AES解密
	origData, err := utils.AESCBCBase64Decode(secretKey, data)

	if err != nil {
		return nil, SECRET_KEY_FAIL, SECRET_KEY_FAIL_MSG
	}

	origLen := len(origData)

	if origLen < DATALEN+TIMELEN+CHANNEL_INFO_LEN || origLen > MAX_DATA {
		return nil, PARAM_FAIL, PARAM_FAIL_MSG
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
		return nil, PARAM_FAIL, PARAM_FAIL_MSG
	}
	if int(reqInfo.DataLen) != origLen {
		return nil, DECODE_FAIL, DECODE_FAIL_MSG
	}
	return reqInfo, SUCCESS, SUCCESS_MSG
}
