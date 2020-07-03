package task

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"tomm/log"
)

type worker struct {
	ID         string
	job        chan *Job
	jobNum     int64
	jobNumLock *sync.RWMutex
	wg         *sync.WaitGroup
	wType      int
}

func newWorker(id string, chanLen int64, wg *sync.WaitGroup) *worker {
	return &worker{
		ID:         id,
		job:        make(chan *Job, chanLen),
		jobNum:     0,
		jobNumLock: &sync.RWMutex{},
		wg:         wg,
	}

}

func (w *worker) close() {
	close(w.job)
}

func (w *worker) doJob(j *Job) bool {
	select {
	case w.job <- j:
		return true
	default:
		return false
	}
}

func (w *worker) startWorker() {

	defer func() {
		if err := recover(); err != nil {
			runtime.Caller(1)
			buf := make([]byte, 1024)
			n := runtime.Stack(buf, false)
			pl := fmt.Sprintf("http server panic: %v\n%s\n", err, buf[:n])
			log.Error("http server recover  is %s", pl)
		}
	}()

	log.Info("Worker Start ID is %s", w.ID)
	for {
		select {

		case job, ok := <-w.job:
			if !ok {
				log.Info("Worker is Closed Id is %s", w.ID)
				w.wg.Done()
				return
			}
			// Do pool
			atomic.AddInt64(&w.jobNum, 1)
			for atomic.CompareAndSwapInt64(&w.jobNum, w.jobNum, w.jobNum+1) {
				log.Debug("Do PoolJob JobID is %d", job.ID)
				res := job.Do()
				if res != nil && job.ResNotify != nil {
					select {
					case job.ResNotify <- res:
					default:
					}
				}
				break
			}

			atomic.AddInt64(&w.jobNum, -1)

			log.Debug("Finish PoolJob JobID is %d", job.ID)
		}
	}
}
