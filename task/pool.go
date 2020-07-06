package task

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
	"tomm/config"
	"tomm/utils"
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
	WorkerContent int64 `yaml:"workerContent"`
}

type Pool struct {
	worker map[string]*worker
	wids   []string
	wg     *sync.WaitGroup
	//wLock  *sync.RWMutex
	conf    *PoolConf
	isClose int32 // 1 表示关闭 2 表示开启
}

func NewPool(conf *PoolConf, wg *sync.WaitGroup) *Pool {

	if conf == nil {
		conf = defaultConf
	}
	p := &Pool{
		wids:    make([]string, 0, conf.WorkerNum),
		worker:  make(map[string]*worker, conf.WorkerNum),
		wg:      wg,
		conf:    conf,
		isClose: 1,
	}

	p.startPool()

	return p
}

func (p *Pool) startPool() {

	for i := 0; i < int(p.conf.WorkerNum); i++ {
		wid, _ := utils.StrUUID()
		w := newWorker(wid, p.conf.WorkerContent, p.wg)
		p.wg.Add(1)
		go w.startWorker()
		p.wids = append(p.wids, wid)
		p.worker[wid] = w
	}

	atomic.AddInt32(&p.isClose, 1)
}

func (p *Pool) DoJob(job *Job) bool {

	if p.isClosed() {
		return false
	}

	w := p.getWork()
	if w == nil {
		return false
	}
	return w.doJob(job)
}

func (p *Pool) getTwoNums() (string, string) {

	rand.Seed(time.Now().UnixNano())
	num1 := rand.Int31n(p.conf.WorkerNum)
	num2 := rand.Int31n(p.conf.WorkerNum)
	for num1 == num2 {
		num2 = rand.Int31n(p.conf.WorkerNum)
	}

	return p.wids[num1], p.wids[num2]
}

func (p *Pool) getWork() *worker {
	// 这里使用p2p 策略来选取 worker
	str1, str2 := p.getTwoNums()
	if str1 == "" || str2 == "" {
		return nil
	}
	w1, ok1 := p.worker[str1]
	w2, ok2 := p.worker[str2]
	if !ok1 || !ok2 {
		return nil
	}
	// TODO: 这里先不考虑 多阶段派任务

	w1.jobNumLock.RLock()
	w2.jobNumLock.RLock()

	defer w1.jobNumLock.RUnlock()
	defer w2.jobNumLock.RUnlock()
	//TODO: 设计上就不支持动态扩容!!!
	// TODO : 这里可以加上job 排队，就是加一个 job的channel 在pool里面多一个go在读这个job channel
	//if w1.jobNum == p.conf.WorkerContent && w2.jobNum == p.conf.WorkerContent {
	//	// 满了则创建新的
	//	w := newWorker()
	//}
	if w1.jobNum == p.conf.WorkerContent && w2.jobNum == p.conf.WorkerContent {
		return nil
	}
	if w1.jobNum < w2.jobNum {
		return w1
	}
	return w2
}

func (p *Pool) Close() {

	if p.isClosed() {
		return
	}

	atomic.AddInt32(&p.isClose, -1)
	//p.cancel()
	for _, v := range p.worker {
		// 这里要在发送端调用 close 才会安全
		v.close()
	}
	p.wg.Wait()

}

func (p *Pool) isClosed() bool {
	if atomic.LoadInt32(&p.isClose) == 1 {
		return true
	}
	return false
}
