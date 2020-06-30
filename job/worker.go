package job

import (
	"context"
	"sync"
)

type worker struct {
	ID         int64
	job        chan *Job
	jobNum     int64
	jobNumLock sync.RWMutex
	wg         *sync.WaitGroup
	ctx        context.Context
}

func newWorker(id int64, chanLen int32, wg *sync.WaitGroup, ctx context.Context) *worker {
	return &worker{
		ID:         id,
		job:        make(chan *Job, chanLen),
		jobNum:     0,
		jobNumLock: sync.RWMutex{},
		wg:         wg,
		ctx:        ctx,
	}

}

func (w *worker) Close() {
	close(w.job)
}

func (w *worker) prepareClose() {
	w.wg.Done()
	// 处理未完成的job
}

func (w *worker) DoJob(job *Job) bool {
	select {
	case w.job <- job:
	default:
		return false
	}
	return true
}

func (w *worker) startWorker() {
	for {
		select {
		case <-w.ctx.Done():
			w.prepareClose()
			return
		case _, ok := <-w.job:
			if !ok {
				w.prepareClose()
				return
			}
			// Do job

		}
	}
}
