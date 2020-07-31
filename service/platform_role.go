package service

import (
	"tomm/api/api"
	"tomm/core/server"
	"tomm/ecode"
	"tomm/log"
	"tomm/service/dao"
)

func AddPlatformRole(c *server.Context) {
	// 这里给id数组就好
	req := api.AddPlatformRoleReq{}

	err := c.Bind(&req)
	if err != nil {
		log.Warn("CheckAdminName Bind Fail err is %s", err.Error())
		server.HttpCode(c, ecode.ParamFail)
	}

	// 检查ids是否正确
	// ....

	err = dao.SavePlatformRole(req.RoleName, req.PlatformIds)

	if err != nil {
		server.HttpCode(c, ecode.SystemFail)
		log.Error("AddPlatformRole SavePlatformRole Err is %s", err.Error())
		return
	}

	server.HttpCode(c, nil)
}

func GetAllPlatformRole(c *server.Context) {
	// page pageSize 分页获取PlatformRole
	// 给一个Total
	req := api.GetAllPlatformRoleReq{}

	err := c.Bind(&req)

	if err != nil {
		log.Error("GetAllPlatformRole Bind Param Err is %s", err.Error())
		server.HttpCode(c, ecode.ParamFail)
		return
	}

	total, err := dao.GetPlatformRoleCount()

	if err != nil {
		log.Error("GetAllPlatformRole GetPlatformRoleCount Err is %s", err.Error())
		server.HttpCode(c, ecode.SystemFail)
		return
	}

	if total == 0 {
		server.HttpData(c, nil)
		return
	}

	infos, err := dao.GetPlatformRoleByPage(req.Page, req.PageSize)

	if err != nil {
		log.Error("GetAllPlatformRole GetPlatformRoleByPage Err is %s", err.Error())
		server.HttpCode(c, ecode.SystemFail)
		return
	}
	res := &api.GetAllPlatformRoleRes{}
	res.Infos = infos
	if len(res.Infos) > 0 {
		res.Total = total
	}
	server.HttpData(c, res)
}

func DeletePlatformRole(c *server.Context) {

	req := api.DeletePlatformRoleReq{}

	err := c.Bind(&req)

	if err != nil {
		log.Error("DeletePlatformRole Bind Param Err is %s", err.Error())
		server.HttpCode(c, ecode.ParamFail)
		return
	}

	err = dao.DeletePlatformRoleByIds(req.Ids)

	if err != nil {
		log.Error("DeletePlatformRole DeletePlatformRoleByIds Err is %s", err.Error())
		server.HttpCode(c, ecode.ParamFail)
		return
	}

	server.HttpCode(c, nil)
}

func UpdatePlatformRole(c *server.Context) {

	req := api.UpdatePlatformRoleReq{}

	err := c.Bind(&req)

	if err != nil {
		log.Error("UpdatePlatformRole Bind Param Err is %s", err.Error())
		server.HttpCode(c, ecode.ParamFail)
		return
	}

	appkeyMap := make(map[string]struct{}, len(req.AppKeys))

	for i := 0; i < len(req.AppKeys); i++ {
		appkeyMap[req.AppKeys[i]] = struct{}{}
	}
	// 查出原来有那些
	err = dao.UpdatePlatformRole(req.RoleSign, req.RoleName, appkeyMap)

	if err != nil {
		log.Error("UpdatePlatformRole UpdatePlatformRole Err is %s", err.Error())
		server.HttpCode(c, ecode.EditFail)
		return
	}

	server.HttpCode(c, nil)
}
