package ecode

import (
	"strconv"
)

var (
	_code = map[int64]struct{}{}
)

type ECodes interface {
	Error() string
	Code() int64
	Equal(ECodes) bool
}

type ECode int64

func addCode(code int64) ECode {
	_, ok := _code[code]

	if ok {
		panic("Cur errMsg is Registered")
	}

	_code[code] = struct{}{}
	return FromInt(code)
}

func (e ECode) Error() string {
	return strconv.FormatInt(int64(e), 10)
}

func (e ECode) Code() int64 { return int64(e) }

func (e ECode) Equal(code ECodes) bool {
	if code == nil {
		code = OK
	}

	return e.Code() == code.Code()
}

func FromInt(code int64) ECode { return ECode(code) }
