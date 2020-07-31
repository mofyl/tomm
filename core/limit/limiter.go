package limit

type Op int

const (
	Success Op = iota
	Drop
)

type DoneInfo struct {
	Err error
	Op  Op
}

type Limiter interface {
	Allow() (func(info DoneInfo), error)
}
