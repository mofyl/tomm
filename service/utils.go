package service

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"tomm/core/server"
	"tomm/ecode"
)

var (
	cli *http.Client
)

func init() {
	cli = &http.Client{}
	cli.Timeout = 5 * time.Second
}

func httpCode(c *server.Context, msgs ecode.ErrMsgs) {
	c.Json(nil, msgs)
}

func httpData(c *server.Context, data interface{}) {
	c.Json(data, ecode.OK)
}

func HttpGet(url string, arg map[string]string) (*http.Response, error) {

	builder := strings.Builder{}

	builder.WriteString(url)

	builder.WriteString("?")
	for k, v := range arg {
		builder.WriteString(fmt.Sprintf("%s=%s&", k, v))
	}

	urlStr := builder.String()
	urlStr = urlStr[:len(urlStr)-1]

	return cli.Get(urlStr)
}
