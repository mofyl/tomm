package task

import (
	"context"
	"fmt"
	"hulk/config"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/sunreaver/logger"
)

var (
	defaultConf = &config.TaskConf{
		TaskNum:       4,
		WorkerNum:     2,
		WorkerContent: 1,
		ExpTime:       60,
	}

	defaultLog logger.Logger
)

type StartTask interface {
	StartTask(ctx *TaskContext) bool
}

type TaskManager struct {
	pool          *Pool
	TaskChan      chan *TaskContext
	doneChan      chan struct{}
	finish        chan struct{}
	wg            *sync.WaitGroup
	resNotifyChan chan *TaskContext
	closeCtx      context.Context
	closeCancel   context.CancelFunc
	isClosed      int32 // 1 表示关闭 2 表示开启
	conf          *config.TaskConf
}

func NewTaskManager(conf *config.TaskConf) *TaskManager {

	if conf == nil {
		conf = defaultConf
	}
	ctx, cancel := context.WithCancel(context.Background())
	tm := &TaskManager{}
	tm.wg = &sync.WaitGroup{}
	tm.TaskChan = make(chan *TaskContext, conf.TaskNum)
	tm.doneChan = make(chan struct{})
	tm.finish = make(chan struct{})
	tm.resNotifyChan = make(chan *TaskContext, conf.TaskNum)
	tm.pool = NewPool(conf, tm.wg)
	tm.conf = conf
	tm.closeCtx = ctx
	tm.closeCancel = cancel
	tm.wg.Add(1)
	go tm.goTask()

	go tm.goNotify()
	tm.isClosed = 2
	return tm
}

func (t *TaskManager) StartTask(ctx *TaskContext) bool {

	if t.isClose() {
		return false
	}

	t.TaskChan <- ctx
	return true
}

func (t *TaskManager) goTask() {

	defer func() {
		if err := recover(); err != nil {
			runtime.Caller(1)
			buf := make([]byte, 1024)
			n := runtime.Stack(buf, false)
			pl := fmt.Sprintf("http server panic: %v\n%s\n", err, buf[:n])
			defaultLog.Errorw("worker recover is ", "recover ", pl)
		}
	}()

	for {
		select {
		case <-t.closeCtx.Done():
			close(t.TaskChan)
			for v := range t.TaskChan {
				t.doJob(v)
			}
			t.pool.Close()

			t.wg.Done()
			return
		case ctx, ok := <-t.TaskChan:
			if !ok {
				break
			}
			t.doJob(ctx)
		}
	}
}

func (t *TaskManager) goNotify() {

	for {
		select {
		case task, ok := <-t.resNotifyChan:
			if !ok {
				break
			}
			t.doNotify(task)
		case _, ok := <-t.doneChan:
			if !ok {
				return
			}
			close(t.resNotifyChan)
			for v := range t.resNotifyChan {
				t.doNotify(v)
			}
			t.finish <- struct{}{}
			return
		}

	}
}

func (t *TaskManager) doJob(ctx *TaskContext) {

	j := &Job{
		ID:        111,
		ResNotify: t.resNotifyChan,
		Do: func() *TaskContext {
			for ctx.curStage < ctx.TaskStage &&
				ctx.TaskHandlers[ctx.curStage](ctx) {
				ctx.curStage++
			}
			return ctx
		},
		IsBlock: ctx.Block,
	}

	t.pool.DoJob(j)
}

func (t *TaskManager) doNotify(ctx *TaskContext) {

	if ctx.NotifyUserChan != nil && ctx != nil {
		if defaultLog != nil {
			defaultLog.Debugw("do Notify")
		}
		ctx.NotifyUserChan <- ctx
	}

}

func (t *TaskManager) Close() {
	if t.isClose() {
		return
	}
	atomic.AddInt32(&t.isClosed, -1)

	t.closeCancel()
	t.wg.Wait()
	t.doneChan <- struct{}{}

	<-t.finish
	close(t.finish)
}

func (t *TaskManager) isClose() bool {
	return atomic.LoadInt32(&t.isClosed) == 1
}
