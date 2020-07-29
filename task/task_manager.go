package task

import (
	"sync"
	"sync/atomic"
)

const (
	MAX_TASK_NUM = 1024
)

type StartTask interface {
	StartTask(ctx *TaskContext) bool
}

type TaskManager struct {
	pool          *Pool
	TaskChan      chan *TaskContext
	doneChan      chan struct{}
	wg            *sync.WaitGroup
	resNotifyChan chan *TaskContext

	isClosed int32 // 1 表示关闭 2 表示开启
}

func NewTaskManager() *TaskManager {
	tm := &TaskManager{}
	tm.wg = &sync.WaitGroup{}
	tm.TaskChan = make(chan *TaskContext, MAX_TASK_NUM)
	tm.doneChan = make(chan struct{})
	tm.resNotifyChan = make(chan *TaskContext, MAX_TASK_NUM)
	tm.pool = NewPool(nil, tm.wg)

	tm.wg.Add(1)
	go tm.doTask()
	tm.isClosed = 2
	return tm
}

//
//func (t *TaskManager) RegisterTaskMD(md *TaskMD, notifyChan chan TaskContext, taskType job.JobApi) *TaskContext {
//
//	if t.isClose() {
//		return nil
//	}
//
//	if md.TaskName == "" {
//		md.TaskName = ""
//	}
//
//	if md.TaskStage <= 0 {
//		panic("Register Task Fail: task stage must be more than zero")
//	}
//
//	if md.TaskHandlers == nil {
//		panic("Register Task Fail: task handlers is nil")
//	}
//
//	t.TaskLock.Lock()
//	_, ok := t.TaskMap[md.TaskName]
//	if ok {
//		t.TaskLock.Unlock()
//		panic("Register Task Fail: task already registered")
//	}
//
//	t.TaskMap[md.TaskName] = md
//	t.TaskLock.Unlock()
//
//	ctx := &TaskContext{
//		TaskName:       md.TaskName,
//		TaskID:         "",
//		Err:            nil,
//		curStage:       0,
//		Type:           taskType,
//		md:             make(map[string]interface{}),
//		NotifyUserChan: notifyChan,
//		st:             t,
//	}
//
//	return ctx
//}

func (t *TaskManager) StartTask(ctx *TaskContext) bool {

	if t.isClose() {
		return false
	}

	select {
	case t.TaskChan <- ctx:
		return true
	default:
		return false
	}
}

func (t *TaskManager) doTask() {

	for {
		select {
		case ctx, ok := <-t.TaskChan:

			if !ok {
				t.doneChan <- struct{}{}
				t.wg.Done()
				return
			}

			t.pool.DoJob(&Job{
				ID:        111,
				ResNotify: t.resNotifyChan,
				Do: func() *TaskContext {
					for ctx.curStage < ctx.TaskStage &&
						ctx.TaskHandlers[ctx.curStage](ctx) {
						ctx.curStage++
					}
					return ctx
				},
			})

		case task, ok := <-t.resNotifyChan:
			if !ok {
				t.wg.Done()
				return
			}

			if task.NotifyUserChan != nil && task != nil {
				select {
				case task.NotifyUserChan <- task:
				default:
				}
			}

		}
	}
}

func (t *TaskManager) Close() {

	if t.isClose() {
		return
	}

	atomic.AddInt32(&t.isClosed, -1)

	close(t.TaskChan)
	<-t.doneChan

	t.pool.Close()
	t.wg.Wait()

	close(t.resNotifyChan)
	close(t.doneChan)
}

func (t *TaskManager) isClose() bool {
	if atomic.LoadInt32(&t.isClosed) == 1 {
		return true
	}
	return false
}
