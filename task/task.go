package task

import (
	"context"
	"unsafe"
)

type JobMarshal interface {
	Marshal() ([]byte, error)
}

type JobType string

type Job struct {
	ID        int64
	ResNotify chan *TaskContext
	Do        func() *TaskContext
}

// 返回false 表示不要进行下一步 true表示要进行下一步
type TaskHandler func(in *TaskContext) bool

// =================================== TaskContext Begin
type TaskContext struct {
	TaskStage    int32 // 从0开始
	TaskHandlers []TaskHandler
	Block bool // 表示是否以阻塞的方式开始任务。若为true表示加不进去的时候就等待
	TaskID   int64
	NotifyUserChan chan *TaskContext
	Err      error

	curStage int32
	md map[string]unsafe.Pointer
	ctx context.Context
	st             StartTask
	createTime int64 // 创建时间
}



func (tc *TaskContext) Set(key string, value unsafe.Pointer) {
	tc.md[key] = value
}

func (tc *TaskContext) Get(key string) (unsafe.Pointer, bool) {
	v, ok := tc.md[key]
	return v, ok
}

func (tc *TaskContext) reset() {
	tc.curStage = 0
	tc.TaskID = 0 // 这里重新生成一个TaskID
	tc.NotifyUserChan = nil
	tc.TaskHandlers = nil
	tc.TaskStage = 0
	tc.ctx = nil
	tc.Err = nil
	for k, _ := range tc.md {
		delete(tc.md, k)
	}
}

func (tc *TaskContext) Start() bool {
	return tc.st.StartTask(tc)
}
