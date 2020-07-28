package service

import (
	"strings"
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

	if !dao.CheckPlatformName(req.PlatformName) {
		log.Debug("RegisterPlatForm Fail Name is Exist")
		server.HttpCode(c, ecode.PlatFormNameFail)
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

	req := api.GetPlatformInfosReq{}

	err := c.Bind(&req)

	if err != nil {
		log.Error("GetPlatformInfos Bind Param Fail Err is %s", err.Error())
		server.HttpCode(c, ecode.ParamFail)
		return
	}
	res := api.GetPlatformInfosRes{}
	count, err := dao.GetPlatformCount()

	if count == 0 {
		server.HttpData(c, res)
		return
	}

	infos, err := dao.GetAllPlatform(req.Page, req.PageSize)

	if err != nil {
		log.Error("GetPlatformInfos Fail Err is %s", err.Error())
		server.HttpCode(c, ecode.SystemFail)
		return
	}

	res.Infos = infos
	if len(res.Infos) > 0 {
		res.Total = count
	}

	server.HttpData(c, res)

}

func DeletePlatform(c *server.Context) {

	req := api.DeletePlatformReq{}

	err := c.Bind(&req)

	if err != nil {
		log.Error("DeletePlatform Bind Param Err is %s", err.Error())
		server.HttpCode(c, ecode.ParamFail)
		return
	}

	err = dao.DeletePlatformByNames(req.Names)

	if err != nil {
		log.Error("DeletePlatform DeletePlatformByIds Err is %s", err.Error())
		server.HttpCode(c, ecode.SystemFail)
		return
	}

	server.HttpCode(c, nil)

}

func GetPlatformByUserID(c *server.Context) {
	// 通过 某个条件 查看 platForm的数据
	req := api.GetPlatformByUserIDReq{}

	err := c.Bind(req)

	if err != nil {
		log.Error("GetAllPlatformRole Err is %s", err.Error())
		server.HttpCode(c, ecode.ParamFail)
		return
	}

	// 查看当前用户的权限组
	userRole, err := dao.GetMMUserPlatformRoleSign(req.UserId)

	if err != nil {
		log.Error("GetPlatformByUserID Fail Err is %s", err.Error())
		server.HttpCode(c, ecode.SystemFail)
		return
	}

	if len(userRole) <= 0 {
		server.HttpCode(c, nil)
		return
	}
	userRoles := strings.Builder{}

	for i := 0; i < len(userRole); i++ {

		userRoles.WriteString(userRole[i].RoleSign)

		if i < len(userRole)-1 {
			userRoles.WriteString(",")
		}

	}

	roleInfo, err := dao.GetPlatformRoleAppKeyByRoleSigns(userRoles.String())

	if err != nil {
		log.Error("GetPlatformRoleAppKeyByRoleSigns Fail Err is %s", err.Error())
		server.HttpCode(c, ecode.SystemFail)
		return
	}

	if len(roleInfo) <= 0 {
		server.HttpCode(c, nil)
		return
	}

	// 合并AppKey
	appkeyMap := make(map[string]struct{}, len(roleInfo))

	for i := 0; i < len(roleInfo); i++ {
		appkeyMap[roleInfo[i].PlatformAppKey] = struct{}{}
	}

	infos := make([]*model.PlatformInfo, 0, len(appkeyMap))
	for k, _ := range appkeyMap {

		info, err := dao.GetPlatformInfo(k)

		if err != nil {
			continue
		}

		infos = append(infos, info)
	}

	res := api.GetPlatformByUserIDRes{
		Infos: infos,
	}

	server.HttpData(c, &res)

}
