package task

import "tomm/ecode"

// 返回false 表示不要进行下一步 true表示要进行下一步
type TaskHandler func(in *TaskContext) bool

type TaskInfo struct {
	TaskName     string // 不同的task name不同
	TaskStep     int32
	TaskHandlers []TaskHandler
	NotifyChan   chan bool
}

//type Task struct {
//	CurStep  int32
//	in       *TaskContext
//	TaskName string // 不同的task name不同
//}

const (
	JobFail    JobType = "JobFail"
	GetUserJob JobType = "GetUserJob"
)

type JobMarshal interface {
	Marshal() ([]byte, error)
}

type JobType string

//
//func NewTaskOut(resType ResType, msgs ecode.ErrMsgs, m JobMarshal) *TaskOut {
//	b, _ := m.Marshal()
//	return &TaskOut{
//		Type: resType,
//		Err:  msgs,
//		Data: b,
//	}
//}

type TaskContext struct {
	TaskName string
	TaskID   string
	Err      ecode.ErrMsgs
	CurStep  int32

	Type       JobType
	Data       []byte
	NotifyChan chan *TaskContext
}

type Job struct {
	ID        int64
	ResNotify chan *TaskContext
	Do        func() *TaskContext
}
