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
}

type errMsg struct {
	ECode
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

func NewErr(err error, code ECode) ErrMsgs {
	return errMsg{
		ECode:  code,
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
