package service

import (
	"time"
	"tomm/api/api"
	"tomm/api/model"
	"tomm/core/server"
	"tomm/ecode"
	"tomm/log"
	"tomm/service/dao"
	"tomm/utils"
)

func VerifyToken(c *server.Context) {
	req := api.VerifyTokenReq{}
	err := c.Bind(&req)
	if err != nil {
		log.Warn("VerifyToken Bind Err is %s", err.Error())
	}

	// 查看token是否存在
	token, expTime, err := dao.GetTokenAndCreate(req.AppKey)

	if err != nil {
		log.Error("Verify Token Fail Err is %s", err.Error())
		server.HttpCode(c, ecode.ParamFail)
		return
	}

	if token != token {
		// 校验失败
		server.HttpCode(c, ecode.VerifyFail)
		return
	}
	res := api.VerifyTokenRes{}
	res.ExpTime = expTime
	server.HttpData(c, res)
}

func GetResourceToken(c *server.Context) {

	req, secretInfo, reqDataInfo, eCode := checkGetTokenReq(c)
	if eCode != nil {
		server.HttpCode(c, eCode)
		return
	}

	// 查看该Code是否存在
	exist, err := dao.MMUserAuthExistDB(model.MMUserAuthInfo{AppKey: req.AppKey})
	if err != nil {
		log.Error("Get Resource Token MMUserAuthExistDB Fail err is %s , Code is %s", err.Error(), reqDataInfo.Code)
		server.HttpCode(c, ecode.CodeFail)
		return
	}

	if !exist {
		log.Error("Get Resource Token MMUserAuthExistDB Code Not Exist,AppKey is %s Code is %s", req.AppKey, reqDataInfo.Code)
		server.HttpCode(c, ecode.CodeFail)
		return
	}

	token, expTime, err := dao.GetTokenAndCreate(req.AppKey)
	if err != nil {
		log.Error("Get Token Fail err is %s", err.Error())
		server.HttpCode(c, ecode.SystemErr)
		return
	}

	tokenInfo := model.TokenInfo{
		Token:      token,
		ExpTime:    expTime,
		ExtendInfo: reqDataInfo.ExtendInfo,
	}
	tokenB, _ := utils.Json.Marshal(tokenInfo)

	resBase64Str, err := utils.AESCBCBase64Encode(secretInfo.SecretKey, tokenB)
	if err != nil {
		log.Error("AESCBCBase64Encode Fail Err is %s", err.Error())
		server.HttpCode(c, ecode.SystemErr)
		return
	}

	res := api.GetTokenRes{
		Token: resBase64Str,
	}

	server.HttpData(c, res)
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

func GetUserInfo(c *server.Context) {
	req := api.GetUserInfoReq{}
	err := c.Bind(&req)

	if err != nil {
		log.Error("Get UserInfo Bind Fail")
		server.HttpCode(c, ecode.ParamFail)
		return
	}
	//token, err := dao.GetToken(req.AppKey)
	//if err != nil {
	//	log.Error("Get UserInfo Token Not Exist Err is %s , AppKey is %s", err.Error(), req.AppKey)
	//	utils.HttpCode(c, ecode.AppKeyFail)
	//	return
	//}
	//
	//if token != req.Token {
	//	log.Error("Get UserInfo Token Not Exist AppKey is %s", req.AppKey)
	//	utils.HttpCode(c, ecode.TokenFail)
	//	return
	//}

	// 使用appkey 获取userID
	// TODO: 这里的Token 先改成userID
	codeInfo, err := dao.GetUserAuthInfo(model.MMUserAuthInfo{AppKey: req.AppKey, MmUserId: req.Token})
	if err != nil {
		log.Error("Get UserInfo Code Info Get Fail Err is %s", err.Error())
		server.HttpCode(c, ecode.ParamFail)
		return
	}

	if codeInfo.Id == 0 {
		log.Error("Get UserInfo Code Info Get Fail")
		server.HttpCode(c, ecode.ParamFail)
		return
	}

	userInfo, errMsg := GetBaseUserInfo(codeInfo.MmUserId)

	if errMsg != nil {
		log.Error("getUserInfo Fail errCode is %d errMsg is %s", errMsg.Code(), errMsg.Error())
		server.HttpCode(c, ecode.MMFail)
		return
	}

	server.HttpData(c, userInfo)

}

func GetUserInfo_V2(c *server.Context) {
	req := api.GetUserInfoReq_V2{}
	err := c.Bind(&req)

	if err != nil {
		log.Error("Get UserInfo Bind Fail")
		server.HttpCode(c, ecode.ParamFail)
		return
	}

	_, err = dao.GetPlatformInfo(req.AppKey)

	if err != nil {
		log.Error("CheckCode GetPlatformInfo Fail err is %s", err.Error())
		server.HttpCode(c, ecode.AppKeyFail)
		return
	}
	// 查看是否授权
	exist, userID, err := dao.CheckCode(req.AppKey, req.Code)
	if err != nil {
		log.Error("checkCode CheckCode Fail err is %s, AppKey is %s , Code is %s", err.Error(), req.AppKey, req.Code)
		server.HttpCode(c, ecode.CodeFail)
		return
	}

	if !exist {
		log.Error("checkCode CheckCode Fail AppKey is %s , Code is %s", req.AppKey, req.Code)
		server.HttpCode(c, ecode.CodeFail)
		return
	}

	// 使用appkey 获取userID
	// TODO: 这里的Token 先改成userID
	codeInfo, err := dao.GetUserAuthInfo(model.MMUserAuthInfo{AppKey: req.AppKey, MmUserId: userID})
	if err != nil {
		log.Error("Get UserInfo Code Info Get Fail Err is %s", err.Error())
		server.HttpCode(c, ecode.ParamFail)
		return
	}

	if codeInfo.Id == 0 {
		log.Error("Get UserInfo Code Info Get Fail")
		server.HttpCode(c, ecode.ParamFail)
		return
	}

	userInfo, errMsg := GetBaseUserInfo(codeInfo.MmUserId)

	if errMsg != nil {
		log.Error("getUserInfo Fail errCode is %d errMsg is %s", errMsg.Code(), errMsg.Error())
		server.HttpCode(c, ecode.MMFail)
		return
	}

	server.HttpData(c, userInfo)

}

func checkGetTokenReq(c *server.Context) (*api.GetTokenReq, *model.PlatformInfo, *api.TokenDataInfo, ecode.ErrMsgs) {
	req := &api.GetTokenReq{}
	err := c.Bind(req)

	if err != nil {
		log.Warn("GetTokenAndCreate Bind Err is %s ", err.Error())
		return nil, nil, nil, ecode.ParamFail
	}

	if req.AppKey == "" || req.Data == "" {
		return nil, nil, nil, ecode.ParamFail
	}

	// 获取该appKey
	secretInfo, err := dao.GetPlatformInfo(req.AppKey)
	if err != nil || secretInfo == nil {
		if err != nil {
			log.Error("GetPlatformInfo Fail AppKey is %s , Err is %s", req.AppKey, err.Error())
		}
		return nil, nil, nil, ecode.AppKeyFail
	}
	//
	//if secretInfo.Id == 0 {
	//	if err != nil {
	//		log.Error("AppKey not illegal AppKey is %s", req.AppKey, err.Error())
	//	}
	//	return nil, nil, nil, ecode.AppKeyFail
	//}

	reqDataInfo, eCode := getTokenDataInfo(secretInfo.SecretKey, req.Data)
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

func getTokenDataInfo(secretKey string, data string) (*api.TokenDataInfo, ecode.ErrMsgs) {

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

	reqInfo := &api.TokenDataInfo{}
	err = utils.Json.Unmarshal(origData, reqInfo)

	if err != nil {
		return nil, ecode.ParamFail
	}
	//if int(reqInfo.DataLen) != origLen {
	//	return nil, DECODE_FAIL, DECODE_FAIL_MSG
	//}
	return reqInfo, nil
}
