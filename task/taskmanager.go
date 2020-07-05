package task

import (
	"sync"
	"tomm/api/job"
)

type StartTask interface{
	StartTask(ctx *TaskContext)bool
}

type TaskManager struct {
	TaskMap       map[string]*TaskMD
	pool          *Pool
	TaskChan      chan *TaskContext
	wg            *sync.WaitGroup
	resNotifyChan chan *TaskContext
}


func NewTaskManager() *TaskManager {
	tm := &TaskManager{}
	tm.TaskMap = make(map[string]*TaskMD)
	tm.wg = &sync.WaitGroup{}
	tm.TaskChan = make(chan *TaskContext, 100)
	tm.resNotifyChan = make(chan *TaskContext , 100)
	tm.pool = NewPool(nil, tm.wg)
	tm.pool.startPool()
	return tm
}

func (t *TaskManager)RegisterTaskMD(md *TaskMD , notifyChan chan *TaskInfo) *TaskContext{
	if md.TaskName == "" {
		md.TaskName = ""
	}

	if md.TaskStage >= 0 {
		panic("Register Task Fail: task stage must be more than zero")
	}

	if md.TaskHandlers == nil {
		panic("Register Task Fail: task handlers is nil")
	}

	_ , ok := t.TaskMap[md.TaskName]
	if ok {
		panic("Register Task Fail: task already registered")
	}

	t.TaskMap[md.TaskName] = md

	ctx := &TaskContext{
		TaskInfo:       &TaskInfo{
			TaskName: md.TaskName,
			TaskID:   "1111",
			Err:      nil,
			CurStage: 0,
			Type:     job.JobApi_JobUserInfo,
			Data:     nil,
		},
		NotifyUserChan: notifyChan,
		st : t ,
	}

	return ctx
}

func (t *TaskManager) StartTask(ctx *TaskContext)bool {

	select{
	case t.TaskChan <- ctx :
		return true
	default:
		return false
	}
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

			t.pool.DoJob(&Job{
				ID:        111,
				ResNotify: t.resNotifyChan,
				Do:        func () *TaskContext{
					for taskInfo.TaskHandlers[task.CurStage](task) {
						task.CurStage++
					}
					return task
				},
			})

		case task, ok := <-t.resNotifyChan:
			if !ok {
				t.wg.Done()
				return
			}

			if task.NotifyUserChan != nil && task.TaskInfo != nil {
				select {
				case task.NotifyUserChan <- task.TaskInfo :
				default:
				}
			}

		}
	}
}

func (t *TaskManager) checkTask(task *TaskContext) (*TaskMD, bool) {

	taskInfo, ok := t.TaskMap[task.TaskName]
	if !ok {
		return nil, false
	}

	if taskInfo.TaskStage < task.CurStage {
		return nil, false
	}

	return taskInfo, true
}

func (t *TaskManager)Close() {
	close(t.TaskChan)
	t.pool.Close()
	close(t.resNotifyChan)
}
