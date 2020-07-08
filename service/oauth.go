package service

import (
	"time"
	"tomm/api/service"
	"tomm/core/server"
	"tomm/ecode"
	"tomm/log"
	"tomm/service/dao"
	"tomm/utils"
)

func (s *Ser) verifyToken(c *server.Context) {
	req := service.VerifyTokenReq{}
	err := c.Bind(&req)
	if err != nil {
		log.Warn("VerifyToken Bind Err is %s", err.Error())
	}

	// 查看token是否存在
	token, expTime, err := dao.GetToken(req.AppKey)

	if err != nil {
		log.Error("Verify Token Fail Err is %s", err.Error())
		httpCode(c, ecode.ParamFail)
	}

	if token != token {
		// 校验失败
		httpCode(c, ecode.VerifyFail)
	}
	res := service.VerifyTokenRes{}
	res.ExpTime = expTime
	httpData(c, res)
}

func (s *Ser) getResourceToken(c *server.Context) {

	req, secretInfo, reqDataInfo, eCode := checkGetTokenReq(c)
	if eCode != nil {
		httpCode(c, eCode)
		return
	}

	// 查看该Code是否存在
	exist, err := dao.CodeExist(service.CodeInfo{AppKey: req.AppKey, Code: reqDataInfo.Code})
	if err != nil {
		httpCode(c, ecode.NewErr(err))
		return
	}

	if exist {
		httpCode(c, ecode.NewErrWithMsg("Code Fail", ecode.RESOURCE_ERR))
		return
	}

	token, expTime, err := dao.GetToken(req.AppKey)
	if err != nil {
		log.Error("Get Token Fail err is %s", err.Error())
		httpCode(c, ecode.NewErr(err))
		return
	}

	tokenInfo := service.TokenInfo{
		Token:      token,
		ExpTime:    expTime,
		ExtendInfo: reqDataInfo.ExtendInfo,
	}
	tokenB, _ := utils.Json.Marshal(tokenInfo)

	resBase64Str, err := utils.AESCBCBase64Encode(secretInfo.SecretKey, tokenB)
	if err != nil {
		log.Error("AESCBCBase64Encode Fail Err is %s", err.Error())
		httpCode(c, ecode.NewErr(err))
		return
	}

	res := service.GetTokenRes{
		Token: resBase64Str,
	}

	httpData(c, res)
	//
	//err, ctx := task.NewTaskContext(s.jobNotify, job.JobApi_JobUserInfo, 3, getTokenJob1, getTokenJob2)
	//
	//if err != nil {
	//	log.Error(err.Error())
	//	return
	//}
	//ctx.Set("secretInfo", secretInfo)
	//ctx.Set("reqDataInfo", reqDataInfo)
	//ctx.Start()

	return
}

func (s *Ser) getUserInfo(c *server.Context) {
	req := service.GetUserInfoReq{}
	err := c.Bind(&req)

	if err != nil {
		httpCode(c, ecode.NewErr(err))
		return
	}

	// 查看Token是否存在
	if !dao.TokenExist(req.AppKey) {
		httpCode(c, ecode.VerifyFail)
		return
	}

	// 使用appkey 获取userID
	codeInfo, err := dao.GetCodeInfo(service.CodeInfo{AppKey: req.AppKey})
	if err != nil {
		httpCode(c, ecode.NewErr(err))
		return
	}

	userInfo, errMsg := GetBaseUserInfo(codeInfo.MmUserId)

	if errMsg != nil {
		httpCode(c, errMsg)
		return
	}

	httpData(c, userInfo)

}

func checkGetTokenReq(c *server.Context) (*service.GetTokenReq, *service.PlatformInfo, *service.TokenDataInfo, ecode.ErrMsgs) {
	req := &service.GetTokenReq{}
	err := c.Bind(req)

	if err != nil {
		log.Warn("GetToken Bind Err is %s ", err.Error())
		return nil, nil, nil, ecode.ParamFail
	}

	if req.AppKey == "" || req.Data == "" {
		return nil, nil, nil, ecode.ParamFail
	}

	// 获取该appKey
	secretInfo, err := dao.GetPlatformInfo(req.AppKey)
	if err != nil || secretInfo == nil {
		if err != nil {
			log.Error("GetToken Fail AppKey is %s , Err is %s", req.AppKey, err.Error())
		}
		return nil, nil, nil, ecode.SecretKeyFail
	}

	reqDataInfo, eCode := GetDataInfo(secretInfo.SecretKey, req.Data)
	if eCode != nil {
		log.Error("Get Data Info Fail ")
		return nil, nil, nil, eCode
	}
	//
	//if reqDataInfo.ChannelInfo == "" ||
	//	reqDataInfo.SendTime == 0 ||
	//	!utils.CheckUrl(reqDataInfo.BackUrl) {
	//	return nil, nil, ecode.ParamFail
	//}
	// 超过10分钟就不处理了
	if time.Now().Unix()-reqDataInfo.SendTime > MAX_TTL {
		log.Error("Package Timeout ")
		return nil, nil, nil, ecode.PackageTimeout
	}

	return req, secretInfo, reqDataInfo, nil
}

func GetDataInfo(secretKey string, data string) (*service.TokenDataInfo, ecode.ErrMsgs) {

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

	reqInfo := &service.TokenDataInfo{}
	err = utils.Json.Unmarshal(origData, reqInfo)

	if err != nil {
		return nil, ecode.ParamFail
	}
	//if int(reqInfo.DataLen) != origLen {
	//	return nil, DECODE_FAIL, DECODE_FAIL_MSG
	//}
	return reqInfo, nil
}
