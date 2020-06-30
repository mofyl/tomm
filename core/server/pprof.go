package server

import (
	"net/http"
	"net/http/pprof"
	"tomm/log"
)

func startPProf(e *Engine) {
	e.GET("/debug/pprof/", pprofHandler(pprof.Index))
	e.GET("/debug/pprof/cmdline", pprofHandler(pprof.Cmdline))
	e.GET("/debug/pprof/profile", pprofHandler(pprof.Profile))
	e.GET("/debug/pprof/symbol", pprofHandler(pprof.Symbol))
	e.GET("/debug/pprof/trace", pprofHandler(pprof.Trace))
	e.wg.Add(1)
	go func() {
		if err := http.ListenAndServe(":9000", nil); err != nil {
			log.Error("pprof error is %s", err.Error())
		}
		e.wg.Done()
	}()
	log.Info("pprof start addr is :9000 ")
}

func pprofHandler(h http.HandlerFunc) HandlerFunc {
	return func(c *Context) {
		h.ServeHTTP(c.Res, c.Req)
	}
}
