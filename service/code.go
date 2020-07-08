package service

import (
	"tomm/api/service"
	"tomm/core/server"
	"tomm/ecode"
	"tomm/service/dao"
	"tomm/utils"
)

func (s *Ser) getCode(c *server.Context) {
	req := service.GetCodeReq{}

	if err := c.Bind(&req); err != nil {
		httpCode(c, ecode.NewErr(err))
		return
	}

	// 检查 code是否存在
	codeInfo, err := dao.GetCodeInfoByUserID(req.UserId)
	if err != nil {
		httpCode(c, ecode.NewErr(err))
		return
	}
	if codeInfo.Id != 0 {
		httpData(c, codeInfo.Code)
		return
	}
	// 检查 app_key 是否存在
	_, err = dao.GetPlatformInfo(req.AppKey)
	if err != nil {
		httpCode(c, ecode.NewErr(err))
		return
	}
	// 创建新的 code
	code, _ := utils.StrUUID()

	codeInfo.AppKey = req.AppKey
	codeInfo.MmUserId = req.UserId
	codeInfo.Code = code
	err = dao.SaveCodeInfo(codeInfo)
	if err != nil {
		httpCode(c, ecode.NewErr(err))
		return
	}

	httpData(c, code)

}
