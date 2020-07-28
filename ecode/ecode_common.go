package ecode

const (
	system_fail   = "System Fail"
	ok            = "Success"
	resource_fail = "resource_fail"
	param_fail    = "Parameter error"
	app_key_fail  = "AppKey error"

	not_value = "not_value"
)

var (
	OK        = NewErrWithMsg(ok, addCode(1))             // 成功
	SystemErr = NewErrWithMsg(system_fail, addCode(-500)) // 服务器错误

	RESOURCE_ERR = addCode(-501) // 服务器对于某个请求 可能某些部分不支持
	//UNKNOW      = addCode(-1000)

	// 1000 ~ 2000  共用的Msg预留
	ParamFail  = NewErrWithMsg(param_fail, addCode(-1000))
	AppKeyFail = NewErrWithMsg(app_key_fail, addCode(-1001))
	SystemFail = NewErrWithMsg(system_fail, addCode(-1002))

	// Redis Error
	NotValue = NewErrWithMsg(param_fail, addCode(-100000))
)
