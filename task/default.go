package task

import (
	"errors"
	"github.com/sunreaver/logger"
	"hulk/config"
	"sync"
	"unsafe"
)

var (
	tm *TaskManager
	ctxPool *sync.Pool
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

func Init(conf config.TaskConf,log logger.Logger) error {


	tm = NewTaskManager(&conf)

	defaultLog = log

	ctxPool = &sync.Pool{
		New: func()interface{} {
			return &TaskContext{
				TaskStage:      0,
				TaskHandlers:   nil,
				TaskID:         0,
				Err:            nil,
				curStage:       0,
				md:             make(map[string]unsafe.Pointer),
				ctx:            nil,
				NotifyUserChan: nil,
				st:             nil,
			}
		},
	}

	return nil
}

func initNoConf(){
	tm = NewTaskManager(nil)
}


func NewTaskContextWithCancel(notifyChan chan *TaskContext , taskStage int32, taskHandlers ...TaskHandler) (*TaskContext ,func () , error) {

	if taskStage <= 0 {
		return nil , nil , errors.New("task stage must more than zero")
	}
	if taskHandlers == nil || len(taskHandlers) <= 0 {
		return nil , nil ,errors.New("task handlers is nil or len is zero")
	}

	ctx := ctxPool.Get().(*TaskContext)
	ctx.TaskStage = taskStage
	ctx.TaskHandlers = taskHandlers
	ctx.TaskID = GetUUID()
	ctx.NotifyUserChan = notifyChan
	ctx.st = tm
	return ctx , func (){
		ctx.reset()
		ctxPool.Put(ctx)
	} , nil
}


func NewTaskContext(notifyChan chan *TaskContext , taskStage int32,isBlock bool, taskHandlers ...TaskHandler) (*TaskContext , error) {

	if taskStage <= 0 {
		return nil , errors.New("task stage must more than zero")
	}
	if taskHandlers == nil || len(taskHandlers) <= 0 {
		return nil ,errors.New("task handlers is nil or len is zero")
	}

	return &TaskContext{
		TaskStage:      taskStage,
		TaskHandlers:   taskHandlers,
		TaskID:         GetUUID(),
		Err:            nil,
		curStage:       0,
		md:             make(map[string]unsafe.Pointer),
		NotifyUserChan: notifyChan,
		st:             tm,
		Block: isBlock,
	} , nil
}




func NewTaskContextWithCtx(notifyChan chan *TaskContext , taskStage int32, v *TaskContext, isBlock bool, taskHandlers ...TaskHandler) (*TaskContext , error) {

	if taskStage <= 0 {
		return nil , errors.New("task stage must more than zero")
	}
	if taskHandlers == nil || len(taskHandlers) <= 0 {
		return nil , errors.New("task handlers is nil or len is zero")
	}
	ctx := &TaskContext{
		TaskStage:      taskStage,
		TaskHandlers:   taskHandlers,
		TaskID:         GetUUID(),
		Err:            nil,
		curStage:       0,
		md:             make(map[string]unsafe.Pointer),
		NotifyUserChan: notifyChan,
		st:             tm,
		Block: isBlock,
	}

	for k , v := range v.md{
		ctx.md[k] =v
	}

	return ctx, nil
}


func Close() {
	tm.Close()
}
