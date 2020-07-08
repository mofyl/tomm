package task

import (
	"tomm/api/job"
	"tomm/ecode"
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

	TaskID   string
	Err      ecode.ErrMsgs
	curStage int32

	Type job.JobApi // 这里表示 使用的是哪个 api号

	md map[string]interface{} // 根据对应的api 来转就好

	NotifyUserChan chan *TaskContext
	st             StartTask
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
	tc.TaskID = "aqweqw" // 这里重新生成一个TaskID
	tc.Err = nil
	for k, _ := range tc.md {
		delete(tc.md, k)
	}
}

func (tc *TaskContext) Start() bool {
	return tc.st.StartTask(tc)
}
