package pool

type Job struct {
	ID        int64
	ResNotify chan []byte
	Do        func() []byte
}
