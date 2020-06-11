package server

import (
	"net/http"
	"time"
)

type EngConfig struct {
	NetWrok      string
	Addr         string
	Timeout      time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type HandlerFunc func(*Context)

type Engine struct {
	cfg    *EngConfig
	Router map[string]map[string][]HandlerFunc // method -> url -> []HandlerFunc

	noRouter []HandlerFunc
}

func NewEngine() *Engine {

}

func (e *Engine) addRouter(method string, path string, handler ...HandlerFunc) {

	if path[0] != '/' {
	}

}

func (e *Engine) SetConfig(cfg *EngConfig) {

}

func (e *Engine) ServeHTTP(r *http.Request, w *http.ResponseWriter) {

}
