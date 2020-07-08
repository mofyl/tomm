package ecode

const (
	system_fail   = "System Fail"
	ok            = "Success"
	resource_fail = "resource_fail"
)

var (
	OK        = NewErrWithMsg(ok, addCode(1))             // 成功
	SystemErr = NewErrWithMsg(system_fail, addCode(-500)) // 服务器错误

	RESOURCE_ERR = addCode(-501) // 服务器对于某个请求 可能某些部分不支持
	//UNKNOW      = addCode(-1000)
)
