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

// =================================== task Begin

// 返回false 表示不要进行下一步 true表示要进行下一步
type TaskHandler func(in *TaskContext) bool

type TaskMD struct {
	TaskName     string // 不同的task name不同
	TaskStage    int32
	TaskHandlers []TaskHandler
	// TODO:这里 Task 也可以有个Type 表示Task的类型 是延时任务 还是实时任务 后面加
}

type TaskInfo struct {
	TaskName string
	TaskID   string
	Err      ecode.ErrMsgs
	CurStage int32

	Type  job.JobApi// 这里表示 使用的是哪个 api号
	Data []byte
}

type TaskContext struct {
	*TaskInfo
	NotifyUserChan chan *TaskInfo
	st StartTask
}

func (tc *TaskContext)Start() bool{
	return tc.st.StartTask(tc)
}
