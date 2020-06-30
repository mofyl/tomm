package ecode

const (
	system_fail = "System Fail"
	ok          = "Success"
)

var (
	OK        = NewErrWithMsg(ok, addCode(1))             // 成功
	SystemErr = NewErrWithMsg(system_fail, addCode(-500)) // 服务器错误

)
