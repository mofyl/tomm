package service

const (

	// 自定义返回码
	SUCCESS     = 200
	SUCCESS_MSG = "Success"

	// 参数错误
	PARAM_FAIL     = 8000
	PARAM_FAIL_MSG = "parameter error"

	// Secret key 获取失败
	SECRET_KEY_FAIL     = 8001
	SECRET_KEY_FAIL_MSG = "Secret Key Can not Find"

	// 解密失败
	DECODE_FAIL     = 8002
	DECODE_FAIL_MSG = "Decode Fail"

	// 数据包超时
	PACKAGE_TIME_OUT     = 8003
	PACKAGE_TIME_OUT_MSG = "Package Time out"

	SYSTEM_FAIL      = 8003
	SYSTEM_FAILE_MSG = "System Fail"

	// TOKEN 校验失败
	VERIFY_FAIL     = 8004
	VERIFY_FAIL_MSG = "Verify Token Fail"
)
