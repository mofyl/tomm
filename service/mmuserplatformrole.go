package service

import (
	"tomm/api/api"
	"tomm/core/server"
	"tomm/ecode"
	"tomm/log"
	"tomm/service/dao"
)

func AddMMPlatformRoles(c *server.Context) {

	req := api.AddMMPlatformRoleReq{}

	err := c.Bind(&req)

	if err != nil {
		log.Error("AddMMPlatformRoles Bind Parameter Fail Err is %s", err.Error())
		server.HttpCode(c, ecode.ParamFail)
		return
	}

	err = dao.AddMMUserPlatformRole(req.MmUserId, req.RoleSigns)

	if err != nil {
		log.Error("AddMMUserPlatformRole Fail Err is %s", err.Error())
		server.HttpCode(c, ecode.SystemFail)
		return
	}

	server.HttpCode(c, nil)

}

func GetMMPlatformRoles(c *server.Context) {

	req := api.GetMMPlatformRolesReq{}

	err := c.Bind(&req)

	if err != nil {
		log.Error("GetMMPlatformRoles Bind Parameter Fail Err is %s", err.Error())
		server.HttpCode(c, ecode.ParamFail)
		return
	}
	// TODO: 这里的交互不知道怎么弄  不知道怎么给
}
