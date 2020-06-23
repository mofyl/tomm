package rending

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"tomm/log"
)

var (
	jsonContentType = []string{"application/json"}
)

type Json struct {
	Data interface{} `json:"data,omitempty"`
}

func (j *Json) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

func (j *Json) Render(w http.ResponseWriter) error {
	writeContentType(w, jsonContentType)
	b, err := json.Marshal(j)
	//b, errmsg := msgpack.Marshal(j)

	if err != nil {
		log.Error("WriteResponse msgPack Marshal ", zap.String("error", err.Error()))
		return err
	}
	if _, err := w.Write(b); err != nil {
		return err
	}
	return nil
}
