package main

import (
	"tomm/service"
)

func main() {
	s := service.NewService()

	s.Start()

	select {}
}
