package server

import (
	"context"
	"net/http"
)

type Context struct {
	r *http.Request
	w http.ResponseWriter

	c context.Context
}
