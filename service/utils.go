package service

import (
	"tomm/core/server"
	"tomm/ecode"
)

func httpCode(c *server.Context, msgs ecode.ErrMsgs) {
	c.Json(nil, msgs)
}

func httpData(c *server.Context, data interface{}) {
	c.Json(data, ecode.OK)
}
