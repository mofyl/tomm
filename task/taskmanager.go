package task

import (
	"sync"
)

type TaskManager struct {
	TaskMap    map[string]*TaskInfo
	pool       *Pool
	TaskChan   chan *TaskContext
	wg         *sync.WaitGroup
	notifyChan chan *TaskContext
}

func NewTaskManager() *TaskManager {
	tm := &TaskManager{}
	tm.TaskMap = make(map[string]*TaskInfo)
	tm.wg = &sync.WaitGroup{}
	tm.TaskChan = make(chan *TaskContext, 100)
	tm.pool = NewPool(nil, tm.wg)
	return tm
}

func (t *TaskManager) doTask() {

	for {
		select {
		case task, ok := <-t.TaskChan:

			if !ok {
				t.wg.Done()
				return
			}

			taskInfo, ok := t.checkTask(task)
			if !ok {
				break
			}

			taskInfo.TaskHandlers[task.CurStep](task)
		case task, ok := <-t.notifyChan:
			if !ok {
				t.wg.Done()
				return
			}

			taskInfo, ok := t.checkTask(task)
			if !ok {
				break
			}

		}
	}
}

func (t *TaskManager) checkTask(task *TaskContext) (*TaskInfo, bool) {

	taskInfo, ok := t.TaskMap[task.TaskName]
	if !ok {
		return nil, false
	}

	if taskInfo.TaskStep < task.CurStep {
		return nil, false
	}

	return taskInfo, true
}
