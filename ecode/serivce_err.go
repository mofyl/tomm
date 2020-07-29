package ecode

const (
	secret_key_fail = "Secret Key Can not Find"

	decode_fail = "Decode Fail"

	package_time_out = "Package Time out"

	verify_fail = "Verify Token Fail"

	code_fail = "Check Code Fail"

	token_fail = "Token Check Fail Please Replace Token"

	vcode_fail         = "Verification Code Fail Please Reflush"
	login_fail         = "Login Name or PassWord Wrong"
	platfrom_name_fail = "Platform Name is Exist"
	edit_fail          = "edit_fail"
	pwd_equal          = "new Pwd Can not equal as old Pwd"
)

var (
	// 8000~ 8999

	SecretKeyFail  = NewErrWithMsg(secret_key_fail, addCode(-8001))
	DeCodeFail     = NewErrWithMsg(decode_fail, addCode(-8002))
	PackageTimeout = NewErrWithMsg(package_time_out, addCode(-8003))
	VerifyFail     = NewErrWithMsg(verify_fail, addCode(-8004))
	CodeFail       = NewErrWithMsg(code_fail, addCode(-8005))
	// TODO 这个Code 表示 Token过期 或非法
	TokenFail = NewErrWithMsg(token_fail, addCode(-8010))

	VCodeFail        = NewErrWithMsg(vcode_fail, addCode(-8011))
	LoginFail        = NewErrWithMsg(login_fail, addCode(-8012))
	PlatFormNameFail = NewErrWithMsg(platfrom_name_fail, addCode(-8013))
	EditFail         = NewErrWithMsg(edit_fail, addCode(-8014))
	PwdEqualFail     = NewErrWithMsg(pwd_equal, addCode(-8015))
)
