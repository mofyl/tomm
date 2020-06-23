package server

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"math"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"tomm/log"
)

const (
	MAX_MEM = 32 << 20 // 32M
)

var (
	default404Body = []byte("404 page not Found")
	default405Body = []byte("405 method not allowed")
	default505Body = []byte("505 System Error")
)

type EngConfig struct {
	Network      string
	Addr         string
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type HandlerFunc func(*Context)

type Engine struct {
	cfg      *EngConfig
	serve    atomic.Value
	wg       *sync.WaitGroup
	Router   map[string]map[string][]HandlerFunc // method -> url -> []HandlerFunc
	Handlers []HandlerFunc
	noRouter []HandlerFunc
	noMethod []HandlerFunc
}

func buildDefaultConf() *EngConfig {
	return &EngConfig{
		Network:      "tcp",
		Addr:         ":8086",
		IdleTimeout:  1 * time.Hour,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
}

func NewEngine(cfg *EngConfig) *Engine {

	if cfg == nil {
		cfg = buildDefaultConf()
	}

	e := &Engine{
		cfg:      cfg,
		wg:       &sync.WaitGroup{},
		Router:   make(map[string]map[string][]HandlerFunc, 4),
		Handlers: make([]HandlerFunc, math.MaxInt8),
		noRouter: make([]HandlerFunc, 0),
		noMethod: make([]HandlerFunc, 0),
	}
	// 加入pprof路由
	startPProf(e)
	return e
}

func (e *Engine) UseFunc(middler ...HandlerFunc) *Engine {
	e.Handlers = append(e.Handlers, middler...)
	return e
}

func (e *Engine) POST(router string, handlers ...HandlerFunc) {
	e.addRouter("POST", router, handlers...)
}

func (e *Engine) GET(router string, handlers ...HandlerFunc) {
	e.addRouter("GET", router, handlers...)
}

func (e *Engine) PUT(router string, handlers ...HandlerFunc) {
	e.addRouter("PUT", router, handlers...)
}

func (e *Engine) DELETE(router string, handlers ...HandlerFunc) {
	e.addRouter("DELETE", router, handlers...)
}
func (e *Engine) addRouter(method string, path string, handler ...HandlerFunc) {

	if path[0] != '/' {
		panic("path mast begin with '/'")
	}

	if method == "" {
		panic("http method can not empty ")
	}

	if len(handler) == 0 {
		panic("must be at least one handler")
	}

	router, _ := e.Router[method]
	if router == nil {
		router = make(map[string][]HandlerFunc)
		e.Router[method] = router
	}

	prelude := func(c *Context) {
		c.Method = method
		c.RouterPath = path
		log.Debug("router ", zap.String("method is ", method), zap.String("path is ", path))
	}

	handlers := append([]HandlerFunc{prelude}, handler...)
	_, ok := router[path]
	if ok {
		panic(fmt.Sprintf("path is exist path is %s , method is %s", path, method))
	}

	router[path] = handlers

	e.NoMethod(func(c *Context) {
		c.Byte(404, "text/plain; chatset=utf-8", default404Body)
		c.Abort()
	})

	e.NoRouter(func(c *Context) {
		c.Byte(405, "text/plain; chatset=utf-8", default405Body)
		c.Abort()
	})
}

func (e *Engine) NoRouter(handler ...HandlerFunc) {
	e.noRouter = handler
	e.rebuild404Handler()
}

func (e *Engine) NoMethod(handler ...HandlerFunc) {
	e.noMethod = handler
	e.rebuild405Handler()
}

func (e *Engine) rebuild404Handler() {
	e.noRouter = e.combineHandlers(e.noRouter...)
}

func (e *Engine) rebuild405Handler() {
	e.noMethod = e.combineHandlers(e.noMethod...)
}

func (e *Engine) combineHandlers(handler ...HandlerFunc) []HandlerFunc {
	finalSize := len(e.Handlers) + len(handler)

	mergeHandler := make([]HandlerFunc, 0, finalSize)

	mergeHandler = append(mergeHandler, e.Handlers...)
	mergeHandler = append(mergeHandler, handler...)
	return mergeHandler
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	c := &Context{
		Res:     w,
		Req:     r,
		index:   -1,
		e:       e,
		handler: nil,
		Ctx:     nil,
	}

	e.HandlerContext(c)
}

func (e *Engine) HandlerContext(ctx *Context) {

	// parseFrom
	cType := ctx.Req.Header.Get("Content-Type")
	switch {
	case strings.Contains(cType, "multipart/form-data"):
		ctx.Req.ParseMultipartForm(MAX_MEM)
	default:
		ctx.Req.ParseForm()
	}
	var cancel context.CancelFunc
	ctx.Ctx, cancel = context.WithCancel(context.TODO())
	defer cancel()

	e.prepareHandler(ctx)
	ctx.Next()
}

func (e *Engine) prepareHandler(c *Context) {
	method := c.Req.Method
	path := c.Req.URL.Path

	routers, ok := e.Router[method]

	if !ok {
		c.handler = e.noMethod
		return
	}

	h, ok := routers[path]

	if !ok {
		c.handler = e.noRouter
		return
	}

	c.handler = h
	return
}

func (e *Engine) RunServer() {
	ser := &http.Server{
		Addr:         e.cfg.Addr,
		Handler:      e,
		IdleTimeout:  e.cfg.IdleTimeout,
		ReadTimeout:  e.cfg.ReadTimeout,
		WriteTimeout: e.cfg.WriteTimeout,
	}

	e.serve.Store(ser)
	e.wg.Add(1)
	go func() {
		if err := ser.ListenAndServe(); err != nil {
			log.Error("RunServer ListenAndServer ", zap.String("error", err.Error()))
		}
		e.wg.Done()
	}()
	log.Info("Http Server Start", zap.String("Addr is ", e.cfg.Addr))
}

func (e *Engine) Close() {
	s := e.serve.Load().(*http.Server)
	s.Close()
	e.wg.Wait()
	log.Info("Http Serve Stop", zap.String("Addr is ", e.cfg.Addr))
}
