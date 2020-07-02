package ecode

const (
	param_fail = "parameter error"

	secret_key_fail = "Secret Key Can not Find"

	decode_fail = "Decode Fail"

	package_time_out = "Package Time out"

	verify_fail = "Verify Token Fail"
)

var (
	// 8000~ 8999
	ParamFail      = NewErrWithMsg(param_fail, addCode(-8000))
	SecretKeyFail  = NewErrWithMsg(secret_key_fail, addCode(-8001))
	DeCodeFail     = NewErrWithMsg(decode_fail, addCode(-8002))
	PackageTimeout = NewErrWithMsg(package_time_out, addCode(-8003))
	VerifyFail     = NewErrWithMsg(verify_fail, addCode(-8004))
)
