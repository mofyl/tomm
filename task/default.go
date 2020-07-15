package task

import (
	"errors"
	"tomm/api/job"
)

var (
	tm *TaskManager
)

//
//func init() {
//
//	defaultConf = &PoolConf{}
//	err := config.Decode(config.CONFIG_FILE_NAME, "pool", defaultConf)
//	if err != nil {
//		panic("Pool Load Config Fail Err is " + err.Error())
//	}
//	tm = NewTaskManager()
//}

func NewTaskContext(notifyChan chan *TaskContext, taskType job.JobApi, taskStage int32, taskHandlers ...TaskHandler) (error, *TaskContext) {

	if taskStage <= 0 {
		return errors.New("task stage must more than zero"), nil
	}
	if taskHandlers == nil || len(taskHandlers) <= 0 {
		return errors.New("task handlers is nil or len is zero"), nil
	}

	return nil, &TaskContext{
		TaskStage:      taskStage,
		TaskHandlers:   taskHandlers,
		TaskID:         "",
		Err:            nil,
		curStage:       0,
		Type:           taskType,
		md:             make(map[string]interface{}),
		NotifyUserChan: notifyChan,
		st:             tm,
	}
}

func Close() {
	tm.Close()
}
