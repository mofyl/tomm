package service

import (
	"tomm/api/service"
	"tomm/core/server"
	"tomm/ecode"
	"tomm/log"
	"tomm/service/dao"
	"tomm/utils"
)

func (s *Ser) getCode(c *server.Context) {
	req := service.GetCodeReq{}

	if err := c.Bind(&req); err != nil {
		httpCode(c, ecode.NewErr(err))
		return
	}

	// 检查 app_key 是否存在
	_, err := dao.GetPlatformInfo(req.AppKey)
	if err != nil {
		log.Error("GetPlatformInfo Fail Err is %s , AppKey is %s", err.Error(), req.AppKey)
		httpCode(c, ecode.AppKeyFail)
		return
	}

	// 检查 code是否存在
	codeInfo, err := dao.GetCodeInfo(service.CodeInfo{MmUserId: req.UserId, AppKey: req.AppKey})
	if err != nil {
		log.Error("GetCodeInfoByUserID Fail Err is %s , UserID is %s", err.Error(), req.UserId)
		httpCode(c, ecode.SystemErr)
		return
	}

	code, _ := utils.StrUUID()
	if codeInfo.Id == 0 {
		// 开始授权
		codeInfo.AppKey = req.AppKey
		codeInfo.MmUserId = req.UserId
		err = dao.SaveCodeInfo(codeInfo, code)
		if err != nil {
			log.Error("SaveCodeInfo Fail Err is %s , Code Info is %v", err.Error(), codeInfo)
			httpCode(c, ecode.SystemErr)
			return
		}
	}

	httpData(c, code)

}

func (s *Ser) checkCode(c *server.Context) {
	req := service.CheckCodeReq{}

	err := c.Bind(&req)

	if err != nil {
		log.Error("checkCode Bind Fail err is %s", err.Error())
		httpCode(c, ecode.ParamFail)
		return
	}

	// 查看是否授权
	exist, userID, err := dao.CheckCode(req.AppKey, req.Code)
	if err != nil {
		log.Error("checkCode CheckCode Fail err is %s, AppKey is %s , Code is %s", err.Error(), req.AppKey, req.Code)
		httpCode(c, ecode.CodeFail)
		return
	}

	if !exist {
		log.Error("checkCode CheckCode Fail AppKey is %s , Code is %s", req.AppKey, req.Code)
		httpCode(c, ecode.CodeFail)
		return
	}

	// 检查成功返回userID
	httpData(c, userID)

}
