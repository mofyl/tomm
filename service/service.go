package service

import (
	"encoding/binary"
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

	return s
}

func (s *Ser) registrRouter() {
	s.e.GET("/getToken")
}

func (s *Ser) getToken(c *server.Context) {
	req := GetTokenReq{}
	err := c.Bind(&req)

	if err != nil {
		log.Warn("GetToken", zap.String("Bind Error ", err.Error()))
		c.Json(200, &BaseMsg{Code: PARAM_FAIL, Msg: PARAM_FAIL_MSG})
		return
	}

	if req.AppKey == "" || req.Data == "" {
		c.Json(200, &BaseMsg{Code: PARAM_FAIL, Msg: PARAM_FAIL_MSG})
		return
	}
	//
	//dataLen := len(req.Data)
	//
	//if dataLen < DATALEN+TIMELEN+CHANNEL_INFO_LEN || dataLen > MAX_DATA {
	//	c.Json(200, &BaseMsg{Code: PARAM_FAIL, Msg: PARAM_FAIL_MSG})
	//	return
	//}

	// 获取该appKey
	secretInfo, err := oauth.GetOAuthInfo(req.AppKey)
	if err != nil || secretInfo == nil {
		if err != nil {
			log.Warn("GetToken", zap.String("Redis Get Key Fail AppKey is ", req.AppKey),
				zap.String("Redis Key Get Fail err is ", err.Error()))
		}
		c.Json(200, &BaseMsg{Code: SECRET_KEY_FAIL, Msg: SECRET_KEY_FAIL_MSG})
		return
	}

	// 使用 secretKey 进行 AES解密
	origData, err := utils.AESCBCBase64Decode(secretInfo.SecretKey, req.Data)

	if err != nil {
		c.Json(200, &BaseMsg{Code: DECODE_FAIL, Msg: DECODE_FAIL_MSG})
		return
	}

	origLen := len(origData)

	if origLen < DATALEN+TIMELEN+CHANNEL_INFO_LEN || origLen > MAX_DATA {
		c.Json(200, &BaseMsg{Code: PARAM_FAIL, Msg: PARAM_FAIL_MSG})
		return
	}

	dataLen := binary.BigEndian.Uint32(origData[:DATALEN])
	sendTime := binary.BigEndian.Uint64(origData[DATALEN : DATALEN+TIMELEN])
	channelInfo := origData[DATALEN+TIMELEN : DATALEN+TIMELEN+CHANNEL_INFO_LEN]
	extendInfo := origData[DATALEN+TIMELEN+CHANNEL_INFO_LEN:]

	if int(dataLen) != origLen {
		c.Json(200, &BaseMsg{Code: DECODE_FAIL, Msg: DECODE_FAIL_MSG})
		return
	}

	// 超过10分钟就不处理了
	if time.Now().Unix()-int64(sendTime) > MAX_TTL {
		c.Json(200, &BaseMsg{Code: PACKAGE_TIME_OUT, Msg: PACKAGE_TIME_OUT_MSG})
		return
	}

	log.Info("Channel Info ", zap.String("Channel ", string(channelInfo)))

	// 构造返回值
	c.Json(200, &BaseMsg{Code: SUCCESS, Msg: SUCCESS_MSG, Data: extendInfo})
}
