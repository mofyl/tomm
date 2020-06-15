package server

import (
	"context"
	"go.uber.org/zap"
	"math"
	"net/http"
	"tomm/core/server/rending"
	"tomm/log"
)

const (
	ABORT_INDEX = math.MaxInt8 / 2
)

type Context struct {
	Req *http.Request
	Res http.ResponseWriter
	Ctx context.Context

	e       *Engine
	index   int8
	handler []HandlerFunc

	Method     string
	RouterPath string
	//Err        message.ErrMsg
}

func (c *Context) Abort() {
	c.index = ABORT_INDEX
}

func (c *Context) Status(code int) {
	c.Res.WriteHeader(code)
}

func (c *Context) Next() {
	c.index++
	for c.handler != nil && c.index < int8(len(c.handler)) {
		c.handler[c.index](c)
		c.index++
	}
}

func (c *Context) Render(code int, render rending.Render) error {
	render.WriteContentType(c.Res)

	if code > 0 {
		c.Status(code)
	}

	err := render.Render(c.Res)
	if err != nil {
		log.Error("Context: Write Response ", zap.String("error", err.Error()))
		return err
	}
	return nil
}

func (c *Context) Json(code int, data interface{}) error {
	return c.Render(code, &rending.Json{
		Data: data,
	})
}

func (c *Context) String(code int, format string, data ...interface{}) error {
	return c.Render(code, &rending.String{
		Format: format,
		Data:   data,
	})
}

func (c *Context) Byte(code int, contentType string, data ...[]byte) error {
	return c.Render(code, &rending.Data{
		ContentType: contentType,
		Data:        data,
	})
}
