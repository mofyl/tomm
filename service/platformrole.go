package service

import (
	"strings"
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

	ids := strings.Split(req.PlatformIds, ",")

	err = dao.SavePlatformRole(req.RoleName, ids)

	if err != nil {
		server.HttpCode(c, ecode.SystemFail)
		log.Error("AddPlatformRole SavePlatformRole Err is %s", err.Error())
		return
	}

	server.HttpCode(c, nil)
}
