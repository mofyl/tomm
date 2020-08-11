package task

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

type worker struct {
	ID     int64
	job    chan *Job
	jobNum int64
	wg     *sync.WaitGroup
	wType  WorkType
	info   *poolInfo

	Blocking uint32 // 1表示正在执行阻塞任务 若该worker正在执行阻塞任务，那么后面便不会给他排任务 直到阻塞任务执行完成
}

func newWorker(id int64, chanLen int64, wg *sync.WaitGroup, info *poolInfo, wType WorkType) *worker {
	return &worker{
		ID:     id,
		job:    make(chan *Job, chanLen),
		jobNum: 0,
		wg:     wg,
		info:   info,
		wType:  wType,
	}

}

func (w *worker) doJob(j *Job) bool {
	w.job <- j
	return true
	//default:
	//	return false
}

func (w *worker) IsBlocking() bool {
	return atomic.LoadUint32(&w.Blocking) == 1
}

func (w *worker) startWorker() {

	defer func() {
		if err := recover(); err != nil {
			runtime.Caller(1)
			buf := make([]byte, 1024)
			n := runtime.Stack(buf, false)
			pl := fmt.Sprintf("http server panic: %v\n%s\n", err, buf[:n])
			defaultLog.Errorw("worker recover is ", "recover ", pl)
		}
	}()

	if defaultLog != nil {
		defaultLog.Debugw("Worker Start ID ", "id is ", w.ID)
	}
	for {
		select {
		case <-w.info.ctx.Done():
			// 将自己从Pool身上删除
			if w.wType == TEMPORARY && w.info.cancel != nil {
				w.info.cancel()
				w.info.f(w.ID)
			}
			close(w.job)
			for v := range w.job {
				doJob(v)
			}
			w.wg.Done()
			return
		case job, ok := <-w.job:
			if !ok {
				//log.Info("Worker is Closed Id is %s", w.ID)
				break
			}
			// Do pool
			//atomic.AddInt64(&w.jobNum, 1)
			atomic.AddInt64(&w.jobNum, 1)
			//log.Debug("Do PoolJob JobID is %d", job.ID)
			if job.IsBlock {
				atomic.StoreUint32(&w.Blocking, 1)
			}
			doJob(job)
			atomic.AddInt64(&w.jobNum, -1)
			atomic.StoreUint32(&w.Blocking, 0)
		}
	}
}

func doJob(job *Job) {
	res := job.Do()
	if res != nil && job.ResNotify != nil {
		job.ResNotify <- res
	}
}
