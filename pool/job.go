package pool

import "tomm/ecode"

type JobRes struct {
	Err  ecode.ErrMsgs
	Data []byte
}

type Job struct {
	ID        int64
	ResNotify chan *JobRes
	Do        func() *JobRes
}
