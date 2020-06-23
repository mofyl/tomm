package service

import "tomm/core/server"

type Ser struct {
	e *server.Engine
}

func NewService() *Ser {
	s := &Ser{}

	e := server.NewEngine(nil)

	s.e = e

	return s
}

func (s *Ser) registrRouter() {

}
