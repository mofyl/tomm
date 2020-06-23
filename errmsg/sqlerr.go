package errmsg

type SqlErr struct {
	ErrMsg
}

func NewSqlErr(msg string) error {
	e := &SqlErr{}

	e.Msg = msg

	return e
}
