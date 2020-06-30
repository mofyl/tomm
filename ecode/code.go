package ecode

import (
	"strconv"
)

var (
	_code = map[int64]struct{}{}
)

type ECodeOp interface {
	SetMsg(msg string)
}

type ECodes interface {
	Error() string
	Code() int64
}

type ECode struct {
	ErrCode int64
	ErrMsg  string
}

func addCode(code int64) ECodes {
	_, ok := _code[code]

	if ok {
		panic("Cur ECode is Registered")
	}

	_code[code] = struct{}{}
	return ECode{
		ErrCode: code,
	}
}

func (e ECode) Code() int64 {
	return e.ErrCode
}

func (e ECode) Error() string {
	if e.ErrMsg != "" {
		return e.ErrMsg
	}
	return strconv.FormatInt(e.ErrCode, 10)
}

func Warp(code ECodes, message string) ECodes {
	if op, ok := code.(ECodeOp); ok {
		op.SetMsg(message)
	}
	return code
}

func (e ECode) SetMsg(msg string) {
	e.ErrMsg = msg
}

type Code int
