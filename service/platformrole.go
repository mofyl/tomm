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

	err := c.Bind(req)
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

	res := &api.GetAllPlatformRoleRes{}
	if total == 0 {
		res.Total = 0
		server.HttpData(c, res)
		return
	}

	infos, err := dao.GetPlatformRoleByPage(req.Page, req.PageSize)

	if err != nil {
		log.Error("GetAllPlatformRole GetPlatformRoleByPage Err is %s", err.Error())
		server.HttpCode(c, ecode.SystemFail)
		return
	}

	res.Infos = infos
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

	// 查出原来有那些

	// 那些需要删除 那些需要增加
	// 是否需要改名字
}
