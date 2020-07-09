package main

import (
	"os"
	"os/signal"
	"syscall"
	"tomm/service"
)

func main() {

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGHUP, syscall.SIGABRT, syscall.SIGQUIT)

	s := service.NewService()

	s.Start()
	<-c

	s.Close()
}
