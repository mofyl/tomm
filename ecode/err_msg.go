package ecode

//
//type errMsg struct {
//	Msg string
//}
//
//func (e *errMsg) Error() string {
//	return e.Msg
//}

type ErrMsgs interface {
	ECodes
	SetMsg(msg string)
	EqualErr(err error) bool
}

type errMsg struct {
	ECode  `json:"err_code"`
	ErrMsg string `json:"err_msg"`
}

func (e errMsg) Error() string {
	if e.ErrMsg != "" {
		return e.ErrMsg
	}
	return e.ECode.Error()
}

func (e errMsg) SetMsg(msg string) {
	e.ErrMsg = msg
}

func (e errMsg) EqualErr(err error) bool {

	if err == nil {
		return false
	}

	return e.ErrMsg == err.Error()
}

func NewErr(err error) ErrMsgs {
	return errMsg{
		ECode:  RESOURCE_ERR,
		ErrMsg: err.Error(),
	}
}

func NewErrWithMsg(msg string, code ECode) ErrMsgs {
	return errMsg{
		ECode:  code,
		ErrMsg: msg,
	}
}

func SetMsgFromErr(err error, code errMsg) ErrMsgs {
	code.SetMsg(err.Error())
	return code
}
