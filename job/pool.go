package job

import (
	"context"
	"math/rand"
	"sync"
	"time"
	"tomm/config"
	"tomm/log"
)

const (
	MAX_WORKER = 10
)

var (
	defaultConf *PoolConf
)

func init() {
	defaultConf = &PoolConf{}
	err := config.Decode(config.CONFIG_FILE_NAME, "pool", defaultConf)
	if err != nil {
		panic("Pool Load Config Fail Err is " + err.Error())
	}
}

type PoolConf struct {
	WorkerNum     int32 `yaml:"workerNum"`
	WorkerContent int32 `yaml:"workerContent"`
}

type Pool struct {
	workerNum int
	worker    []*worker
	wg        *sync.WaitGroup
	wLock     *sync.RWMutex
	conf      *PoolConf
	ctx       context.Context
	cancel    context.CancelFunc
}

func newPool(conf *PoolConf) *Pool {
	if conf == nil {
		conf = defaultConf
	}
	ctx, cancel := context.WithCancel(context.TODO())
	p := &Pool{}
	p.wg = &sync.WaitGroup{}
	p.wLock = &sync.RWMutex{}
	p.conf = conf
	p.ctx = ctx
	p.cancel = cancel
	return p
}

func (p *Pool) StartPool() {

	for i := 0; i < int(p.conf.WorkerNum); i++ {
		id := int64(i + 1)
		w := newWorker(id, p.conf.WorkerContent, p.wg, p.ctx)
		p.wg.Add(1)
		go w.startWorker()
		p.worker = append(p.worker, w)
		log.Info("Worker Start ID is %d", id)
	}

}

func (p *Pool) DoJob(job *Job) bool {

	w := p.GetWork()

	return w.DoJob(job)
}

func (p *Pool) getTwoNums() (int32, int32) {
	rand.Seed(time.Now().Unix())
	num1 := rand.Int31n(p.conf.WorkerNum)
	num2 := rand.Int31n(p.conf.WorkerNum)
	if num1 == num2 {
		for num1 != num2 {
			num2 = rand.Int31n(p.conf.WorkerNum)
		}
	}
	return num1, num2
}

func (p *Pool) GetWork() *worker {
	// 这里使用p2p 策略来选取 worker
	num1, num2 := p.getTwoNums()

	w1 := p.worker[num1]
	w2 := p.worker[num2]

	return nil
}
