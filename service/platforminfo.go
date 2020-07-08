package service

import (
	"tomm/api/service"
	"tomm/core/server"
	"tomm/ecode"
	"tomm/service/dao"
)

func (s *Ser) registerPlatform(c *server.Context) {
	req := service.RegisterPlatformReq{}
	err := c.Bind(&req)

	if err != nil {
		httpCode(c, ecode.NewErr(err))
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
		httpCode(c, ecode.NewErr(err))
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
		httpCode(c, ecode.NewErr(err))
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
