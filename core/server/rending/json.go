package rending

import (
	"encoding/json"
	"net/http"
	"tomm/log"
)

var (
	jsonContentType = []string{"application/json"}
)

type Json struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"json,omitempty"`
}

func (j *Json) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

func (j *Json) Render(w http.ResponseWriter) error {
	writeContentType(w, jsonContentType)

	b, err := json.Marshal(j)
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
