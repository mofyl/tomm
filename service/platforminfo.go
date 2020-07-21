package service

import (
	"tomm/api/api"
	"tomm/api/model"
	"tomm/core/server"
	"tomm/ecode"
	"tomm/log"
	"tomm/service/dao"
)

func RegisterPlatform(c *server.Context) {
	req := api.RegisterPlatformReq{}
	err := c.Bind(&req)

	if err != nil {
		log.Error("RegisterPlatForm Fail Parameter Wrong Err is %s", err.Error())
		server.HttpCode(c, ecode.ParamFail)
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
		server.HttpCode(c, ecode.SystemErr)
		return
	}

	res := api.RegisterPlatformRes{
		SecretKey: info.SecretKey,
		AppKey:    info.AppKey,
	}

	server.HttpData(c, &res)

}

func CheckPlatformName(c *server.Context) {
	req := api.CheckPlatformNameReq{}
	err := c.Bind(&req)

	if err != nil {
		log.Error("checkPlatformName Bind Parameter Fail Err is %s", err.Error())
		server.HttpCode(c, ecode.ParamFail)
		return
	}
	canUsed := dao.CheckPlatformName(req.Name)

	res := &api.CheckPlatformNameRes{}

	if canUsed {
		res.Res = 2
	} else {
		res.Res = 1
	}

	server.HttpData(c, res)
}

func GetPlatformInfos(c *server.Context) {

	infos, err := dao.GetAllPlatform()

	if err != nil {
		log.Error("GetPlatformInfos Fail Err is %s", err.Error())
		server.HttpCode(c, ecode.SystemFail)
		return
	}

	server.HttpData(c, infos)

}

func deletePlatformName(c *server.Context) {

}

func getPlatformByUserID(c *server.Context) {
	// 通过 某个条件 查看 platForm的数据

	// 查询当前用户的权限组

}
