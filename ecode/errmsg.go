package ecode

type ErrMsg struct {
	Msg string
}

func (e *ErrMsg) Error() string {
	return e.Msg
}
