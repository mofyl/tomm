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

// =================================== TaskMD Begin
type TaskMD struct {
	TaskName     string // 不同的task name不同
	TaskStage    int32  // 从0开始
	TaskHandlers []TaskHandler
	// TODO:这里 Task 也可以有个Type 表示Task的类型 是延时任务 还是实时任务 后面加
}

func NewTaskMD(name string, handler ...TaskHandler) *TaskMD {
	if handler == nil {
		panic("Task MD handler is nil")
	}

	lenH := int32(len(handler))

	if lenH <= 0 {
		panic("Task MD handler len must be more than zero")
	}

	return &TaskMD{
		TaskName:     name,
		TaskStage:    int32(len(handler)),
		TaskHandlers: handler,
	}
}

// =================================== TaskContext Begin
type TaskContext struct {
	TaskName string
	TaskID   string
	Err      ecode.ErrMsgs
	CurStage int32

	Type job.JobApi // 这里表示 使用的是哪个 api号

	md map[string]interface{} // 根据对应的api 来转就好

	NotifyUserChan chan TaskContext
	st             StartTask
}

func (tc *TaskContext) Set(key string, value interface{}) {
	tc.md[key] = value
}

func (tc *TaskContext) Get(key string) (interface{}, bool) {
	v, ok := tc.md[key]
	return v, ok
}

func (tc *TaskContext) Start() bool {
	return tc.st.StartTask(tc)
}
