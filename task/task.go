package task

import (
	"context"
)

type JobMarshal interface {
	Marshal() ([]byte, error)
}

type JobType string

type Job struct {
	ID        int64
	ResNotify chan *TaskContext
	Do        func() *TaskContext
	IsBlock   bool
}

// 返回false 表示不要进行下一步 true表示要进行下一步
type TaskHandler func(in *TaskContext) bool

// =================================== TaskContext Begin
type TaskContext struct {
	TaskStage      int32 // 从0开始
	TaskHandlers   []TaskHandler
	Block          bool // 表示该任务是否要阻塞，若是阻塞任务则会影响调度 一个worker 若正在执行阻塞任务，那么后面就不会给这个worker派任务
	TaskID         int64
	NotifyUserChan chan *TaskContext
	Err            error

	curStage int32
	md       map[string]interface{}
	ctx      context.Context
	st       StartTask
	//createTime int64 // 创建时间
}

func (tc *TaskContext) Set(key string, value interface{}) {
	tc.md[key] = value

}

func (tc *TaskContext) Get(key string) (interface{}, bool) {
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
	for k := range tc.md {
		delete(tc.md, k)
	}
}

func (tc *TaskContext) Start() bool {
	return tc.st.StartTask(tc)
}
