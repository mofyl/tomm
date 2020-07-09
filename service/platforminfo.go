package service

import (
	"tomm/api/service"
	"tomm/core/server"
	"tomm/ecode"
	"tomm/log"
	"tomm/service/dao"
)

func (s *Ser) registerPlatform(c *server.Context) {
	req := service.RegisterPlatformReq{}
	err := c.Bind(&req)

	if err != nil {
		log.Error("RegisterPlatForm Fail Parameter Wrong Err is %s", err.Error())
		httpCode(c, ecode.ParamFail)
		return
	}

	info := service.PlatformInfo{
		Memo:        req.Memo,
		IndexUrl:    req.IndexUrl,
		ChannelName: req.PlatformName,
		SignUrl:     req.SignUrl,
	}

	err = dao.CreateOAuthInfo(&info)

	if err != nil {
		log.Error("RegisterPlatForm Create OAuthInfo Fail Err is %s , info is %v", err.Error(), info)
		httpCode(c, ecode.SystemErr)
		return
	}

	res := service.RegisterPlatformRes{
		SecretKey: info.SecretKey,
		AppKey:    info.AppKey,
	}

	httpData(c, &res)

}

func (s *Ser) checkPlatformName(c *server.Context) {
	req := service.CheckPlatformNameReq{}
	err := c.Bind(&req)

	if err != nil {
		log.Error("checkPlatformName Bind Parameter Fail Err is %s", err.Error())
		httpCode(c, ecode.ParamFail)
		return
	}
	canUsed := dao.CheckPlatformName(req.Name)

	res := &service.CheckPlatformNameRes{}

	if canUsed {
		res.Res = 2
	} else {
		res.Res = 1
	}

	httpData(c, res)
}

func (s *Ser) deletePlatformName(c *server.Context) {

}
