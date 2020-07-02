package rending

import (
	"net/http"
	"tomm/log"
	"tomm/utils"
)

var (
	jsonContentType = []string{"application/json"}
)

type JsonEncode interface {
	ToJson() []byte
}

type Json struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
	Data JsonEncode
}

func (j *Json) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

func (j *Json) Render(w http.ResponseWriter) error {
	writeContentType(w, jsonContentType)

	b, err := utils.Json.Marshal(j)
	//b, ecode := msgpack.Marshal(j)

	if err != nil {
		log.Error("WriteResponse msgPack Marshal Err is %s ", err.Error())
		return err
	}
	if _, err := w.Write(b); err != nil {
		return err
	}

	return nil
}
