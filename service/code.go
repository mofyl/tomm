package service

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"
	"tomm/api/api"
	"tomm/api/model"
	"tomm/core/server"
	"tomm/ecode"
	"tomm/log"
	"tomm/redis"
	"tomm/service/dao"
	"tomm/utils"
)

const (
	CODE_DATA_LEN = 12
	CODE_TIME_LEN = 8
	CODE_EXP_TIME = 180 // 3*60s  3min
)

func GetCode(c *server.Context) {
	req := api.GetCodeReq{}

	if err := c.Bind(&req); err != nil {
		server.HttpCode(c, ecode.NewErr(err))
		return
	}
	// 检查 该用户是否存在
	_, errMsg := GetBaseUserInfo(req.UserId)

	if errMsg != nil {
		log.Error("GetCode Check Fail errCode is %d errMsg is %s", errMsg.Code(), errMsg.Error())
		server.HttpCode(c, ecode.MMFail)
		return
	}

	// 检查 app_key 是否存在
	platFormInfo, err := dao.GetPlatformInfo(req.AppKey)
	if err != nil {
		log.Error("GetPlatformInfo Fail Err is %s , AppKey is %s", err.Error(), req.AppKey)
		server.HttpCode(c, ecode.AppKeyFail)
		return
	}

	// 检查 code是否存在
	codeInfo, err := dao.GetUserAuthInfo(model.MMUserAuthInfo{MmUserId: req.UserId, AppKey: req.AppKey})
	if err != nil {
		log.Error("GetUserAuthInfoByUserID Fail Err is %s , UserID is %s", err.Error(), req.UserId)
		server.HttpCode(c, ecode.SystemErr)
		return
	}

	code, _ := utils.StrUUID()
	if codeInfo.Id == 0 {
		// 开始授权
		codeInfo.AppKey = req.AppKey
		codeInfo.MmUserId = req.UserId
		err = dao.SaveMMUserAuthInfo(codeInfo)
		if err != nil {
			log.Error("SaveMMUserAuthInfo Fail Err is %s , Code Info is %v", err.Error(), codeInfo)
			server.HttpCode(c, ecode.SystemErr)
			return
		}
	}

	// 将Code保存到redis
	err = redis.Set(context.TODO(), fmt.Sprintf(dao.CODE_KEY, codeInfo.AppKey, code), codeInfo.MmUserId, dao.CODE_EXP)

	res := api.GetCodeRes{
		Code:    code,
		BackUrl: platFormInfo.SignUrl,
	}
	server.HttpData(c, res)
}

func CheckCode(c *server.Context) {
	req := api.CheckCodeReq{}

	err := c.Bind(&req)

	if err != nil {
		log.Error("checkCode Bind Fail err is %s", err.Error())
		server.HttpCode(c, ecode.ParamFail)
		return
	}
	platformInfo, err := dao.GetPlatformInfo(req.AppKey)

	if err != nil {
		log.Error("CheckCode GetPlatformInfo Fail err is %s", err.Error())
		server.HttpCode(c, ecode.AppKeyFail)
		return
	}

	// TimeStamp+Code
	dataInfo, errCode := getCheckDataInfo(platformInfo.SecretKey, req.Data)
	if errCode != nil {
		server.HttpCode(c, errCode)
		return
	}
	// 查看是否授权
	exist, userID, err := dao.CheckCode(req.AppKey, dataInfo.Code)
	if err != nil {
		log.Error("checkCode CheckCode Fail err is %s, AppKey is %s , Code is %s", err.Error(), req.AppKey, dataInfo.Code)
		server.HttpCode(c, ecode.CodeFail)
		return
	}

	if !exist {
		log.Error("checkCode CheckCode Fail AppKey is %s , Code is %s", req.AppKey, dataInfo.Code)
		server.HttpCode(c, ecode.CodeFail)
		return
	}

	res := api.CheckCodeRes{
		UserId: userID,
	}
	// 检查成功返回userID
	server.HttpData(c, res)
}

func getCheckDataInfo(secretKey string, data string) (model.CheckCodeData, ecode.ErrMsgs) {

	res := model.CheckCodeData{}
	oriData, err := utils.AESCBCBase64Decode(secretKey, data)

	if err != nil {
		log.Error("GetDataInfo data AESCBCBase64Decode Fail err is %s", err.Error())
		return res, ecode.DeCodeFail
	}

	if oriData == nil || len(oriData) < CODE_DATA_LEN {
		log.Error("GetDataInfo data too small")
		return res, ecode.DeCodeFail
	}

	timeStamp := int64(binary.BigEndian.Uint64(oriData[:CODE_TIME_LEN]))

	nowTime := time.Now().Unix()

	if nowTime-timeStamp > CODE_EXP_TIME {
		log.Error("GetDataInfo timeStamp Wrong")
		return res, ecode.DeCodeFail
	}

	//code :=
	res.TimeStamp = timeStamp
	res.Code = string(oriData[CODE_TIME_LEN:])
	return res, nil
}
