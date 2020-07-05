package task

import "testing"

//
func TeskTask(t *testing.T) {
	notify := make(chan *TaskContext)
	tm := NewTaskManager()
	md := TaskMD{
		TaskName:     "111111",
		TaskStage:    1,
		TaskHandlers: ,
		NotifyChan:   nil,
	}
}
