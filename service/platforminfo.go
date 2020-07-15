package service

import (
	"tomm/api/api"
	"tomm/api/model"
	"tomm/core/server"
	"tomm/ecode"
	"tomm/log"
	"tomm/service/dao"
)

func (s *Ser) registerPlatform(c *server.Context) {
	req := api.RegisterPlatformReq{}
	err := c.Bind(&req)

	if err != nil {
		log.Error("RegisterPlatForm Fail Parameter Wrong Err is %s", err.Error())
		httpCode(c, ecode.ParamFail)
		return
	}

	info := model.PlatformInfo{
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

	res := api.RegisterPlatformRes{
		SecretKey: info.SecretKey,
		AppKey:    info.AppKey,
	}

	httpData(c, &res)

}

func (s *Ser) checkPlatformName(c *server.Context) {
	req := api.CheckPlatformNameReq{}
	err := c.Bind(&req)

	if err != nil {
		log.Error("checkPlatformName Bind Parameter Fail Err is %s", err.Error())
		httpCode(c, ecode.ParamFail)
		return
	}
	canUsed := dao.CheckPlatformName(req.Name)

	res := &api.CheckPlatformNameRes{}

	if canUsed {
		res.Res = 2
	} else {
		res.Res = 1
	}

	httpData(c, res)
}

func (s *Ser) deletePlatformName(c *server.Context) {

}

func (s *Ser) getPlatformByUserID(c *server.Context) {

	// 将用户的UserID传过来

}
