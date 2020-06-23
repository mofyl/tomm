package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

type JsonStrcut struct {
	A string `form:"a_str"`
	B string `form:"b_str"`
}

func TestEngine(t *testing.T) {

	c := make(chan os.Signal)
	signal.Notify(c, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP)

	js := JsonStrcut{}
	e := NewEngine(nil)
	go func() {
		<-c
		fmt.Println("signal come")
		e.Close()
	}()
	e.GET("/binding", func(c *Context) {
		if err := c.Bind(&js); err != nil {
			c.String(200, "%s", err.Error())
		}
		c.String(200, "%s", "helloworld"+js.A)
	})

	e.RunServer()

	select {}
}

func TestServer(t *testing.T) {

	c := make(chan os.Signal)
	signal.Notify(c, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP)

	js := JsonStrcut{
		A: "qweqw",
		B: "zczxc",
	}
	e := NewEngine(nil)
	go func() {
		<-c
		fmt.Println("signal come")
		e.Close()
	}()
	e.GET("/", func(c *Context) {
		c.String(200, "%s", "helloworld")
	})
	e.GET("/test1", func(c *Context) {
		c.Json(200, js)
	})
	e.GET("/test2", func(c *Context) {
		c.Byte(200, "application/json; chatset=utf-8", []byte("test2"))
	})

	e.RunServer()

	select {}

}
