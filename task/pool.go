package task

import (
	"context"
	"fmt"
	"hulk/config"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type WorkType int32

var (
	ETERNAL   WorkType = 1
	TEMPORARY WorkType = 2
)

type RemoveFromParent func(id int64)

type poolInfo struct {
	cancel context.CancelFunc
	ctx    context.Context
	f      RemoveFromParent
}

type Pool struct {
	eternalWorker []*worker
	cancelCtx     context.Context
	cancel        context.CancelFunc
	wg            *sync.WaitGroup
	conf          *config.TaskConf
	isClose       int32 // 1 表示关闭 2 表示开启
	randSource    *rand.Rand

	temporaryWorker []*worker
	lockTemporary   *sync.Mutex
}

func NewPool(conf *config.TaskConf, wg *sync.WaitGroup) *Pool {

	if conf == nil {
		conf = defaultConf
	}

	if conf.WorkerNum < 2 {
		panic("NewPool: Worker Num at least 2")
	}

	ctx, cancel := context.WithCancel(context.Background())

	p := &Pool{
		eternalWorker: make([]*worker, 0, conf.WorkerNum),
		wg:            wg,
		conf:          conf,
		isClose:       1,
		randSource:    rand.New(rand.NewSource(time.Now().UnixNano())),
		cancelCtx:     ctx,
		cancel:        cancel,
		lockTemporary: &sync.Mutex{},
	}

	p.startPool()

	return p
}

func (p *Pool) startPool() {

	info := &poolInfo{
		cancel: p.cancel,
		ctx:    p.cancelCtx,
	}

	for i := 0; i < int(p.conf.WorkerNum); i++ {
		wid := i + 1

		w := newWorker(int64(wid), p.conf.WorkerContent, p.wg, info, ETERNAL)
		p.wg.Add(1)
		go w.startWorker()
		//p.wids = append(p.wids, int64(wid))
		p.eternalWorker = append(p.eternalWorker, w)
		//p.worker[int64(i)] = w
	}

	atomic.AddInt32(&p.isClose, 1)
}

func (p *Pool) DoJob(job *Job) bool {

	if p.isClosed() {
		//defaultLog.Debugw("Worker is Closed")
		return false
	}

	w := p.getWorkFormEnternal()

	// if w == nil || w.IsBlocking() {
	// 	w = p.getWorkStep(3, 1, p.getWorkFormEnternal)
	// }

	// if w == nil || w.IsBlocking() {
	// 	w = p.getWorkStep(3, 1, p.getWorkerFromTemp)
	// }

	// if w == nil || w.IsBlocking() {
	// 	w = p.newTempWorker()
	// }

	if w == nil {
		w = p.getWorkStep(3, 1, p.getWorkFormEnternal)
	}

	if w == nil {
		w = p.getWorkStep(3, 1, p.getWorkerFromTemp)
	}
	if w == nil {
		w = p.newTempWorker()
	}

	return w.doJob(job)
}

func (p *Pool) getWorkStep(num int, duration int32, f func() *worker) *worker {

	ticker := time.NewTicker(time.Duration(duration) * time.Second)

	n := 0
	var w *worker
	for n < num {
		<-ticker.C
		n++
		w = f()
		if w == nil {
			continue
		}
		if w.IsBlocking() {
			continue
		}

	}
	ticker.Stop()

	return w

}

func (p *Pool) getTwoNums(num int) (int, int) {

	p.randSource.Seed(time.Now().UnixNano())
	num1 := p.randSource.Intn(num)
	num2 := p.randSource.Intn(num)
	for num1 == num2 {
		num2 = rand.Intn(num)
	}

	return num1, num2
}

func (p *Pool) getTwoWorker(num int, workers []*worker) (*worker, *worker) {
	// 这里使用p2c 策略来选取 worker
	num1, num2 := p.getTwoNums(num)
	return workers[num1], workers[num2]
}

func (p *Pool) getWorkerFromTemp() *worker {

	if p.temporaryWorker == nil {
		p.temporaryWorker = make([]*worker, 0, p.conf.WorkerNum)
		return nil
	}

	if len(p.temporaryWorker) == 0 {
		return nil
	}
	p.lockTemporary.Lock()
	defer p.lockTemporary.Unlock()
	if len(p.temporaryWorker) == 1 {
		w := p.temporaryWorker[0]
		// if atomic.CompareAndSwapInt64(&w.jobNum, p.conf.WorkerContent, w.jobNum) || w.IsBlocking() {
		// 	return nil
		// } else {
		// 	return w
		// }

		if atomic.CompareAndSwapInt64(&w.jobNum, p.conf.WorkerContent, w.jobNum) {
			return nil
		} else {
			return w
		}
	}

	w1, w2 := p.getTwoWorker(len(p.temporaryWorker), p.temporaryWorker)
	if atomic.CompareAndSwapInt64(&w1.jobNum, p.conf.WorkerContent, w1.jobNum) &&
		atomic.CompareAndSwapInt64(&w2.jobNum, p.conf.WorkerContent, w2.jobNum) {
		return nil
	}

	if atomic.LoadInt64(&w1.jobNum) < atomic.LoadInt64(&w2.jobNum) {
		return w1
	} else {
		return w2
	}

}

func (p *Pool) newTempWorker() *worker {

	ctx, cancel := context.WithTimeout(p.cancelCtx, time.Duration(p.conf.ExpTime)*time.Second)
	t := &poolInfo{
		cancel: cancel,
		ctx:    ctx,
		f:      p.removeFromParent,
	}

	w := newWorker(GetUUID(), p.conf.WorkerContent, p.wg, t, TEMPORARY)
	p.wg.Add(1)
	go w.startWorker()
	p.lockTemporary.Lock()
	p.temporaryWorker = append(p.temporaryWorker, w)
	p.lockTemporary.Unlock()
	return w
}

func (p *Pool) removeFromParent(id int64) {

	if p.temporaryWorker == nil || len(p.temporaryWorker) == 0 {
		return
	}
	p.lockTemporary.Lock()
	defer p.lockTemporary.Unlock()

	lenTemp := len(p.temporaryWorker)

	if lenTemp == 1 {
		if p.temporaryWorker[0].ID != id {
			return
		} else {
			p.temporaryWorker = p.temporaryWorker[:0]
		}
	}

	l := 0
	r := lenTemp - 1

	for l < r {
		mid := l + ((l - r) / 2)
		if p.temporaryWorker[mid].ID == id {
			p.temporaryWorker = append(p.temporaryWorker[:mid], p.temporaryWorker[mid+1:]...)
			return
		}
		if p.temporaryWorker[mid].ID > id {
			r = mid - 1
		} else if p.temporaryWorker[mid].ID < id {
			l = mid + 1
		}
	}

}

func (p *Pool) getWorkFormEnternal() *worker {
	// 这里使用p2c 策略来选取 worker
	w1, w2 := p.getTwoWorker(int(p.conf.WorkerNum), p.eternalWorker)

	if defaultLog != nil {
		defaultLog.Debugw(fmt.Sprintf("Cur Select Worker Num is %d , jobNum is %d , Num is %d , jobNum is %d",
			w1.ID, atomic.LoadInt64(&w1.jobNum), w2.ID, atomic.LoadInt64(&w2.jobNum)))
	}

	if atomic.CompareAndSwapInt64(&w1.jobNum, p.conf.WorkerContent, w1.jobNum) &&
		atomic.CompareAndSwapInt64(&w2.jobNum, p.conf.WorkerContent, w2.jobNum) {
		return nil
	}

	if atomic.LoadInt64(&w1.jobNum) < atomic.LoadInt64(&w2.jobNum) {
		return w1
	}

	return w2
}

func (p *Pool) Close() {
	fmt.Println("Close Come")
	if p.isClosed() {
		return
	}

	atomic.AddInt32(&p.isClose, -1)
	//p.cancel()
	p.cancel()

	fmt.Println("Closed")
}

func (p *Pool) isClosed() bool {
	return atomic.LoadInt32(&p.isClose) == 1
}
