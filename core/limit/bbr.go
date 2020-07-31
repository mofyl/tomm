package limit

import (
	"math"
	"sync/atomic"
	"time"
	"tomm/core/metric"
)

var (
	initTime    = time.Now()
	defaultConf = &BBRConfig{
		Enable:       true,
		Window:       time.Second * 10,
		WinBucket:    100,
		CPUThreshold: 800, // 8核心的情况下
	}
)

type cpuGetter func() int64

type BBR struct {
	cpu             cpuGetter              // 获取当前CPU的使用率
	passStat        *metric.RollingCounter // 通过的请求 的 采样窗口
	rtStat          *metric.RollingCounter // 通过的请求RTT的 采用窗口
	inFlight        int64                  // 当前在处理的 请求 数量
	winBucketPerSec int64                  // 1s 里面有多少个桶 也就是1s的采样次数
	prevDrop        atomic.Value           // 上一次 丢弃请求的时间  也就是 上次丢弃请求的时间 距离 initTime经过的时间
	rawMaxPass      int64                  // 上一次的最大通过率
	rawMinRt        int64                  // 上一次的最小 Rtt
	conf            *BBRConfig
}

// 比如要设定 10s内采样100次 那么Window就是10s  WinBucket就是100
type BBRConfig struct {
	Enable       bool          // 是否开启
	Window       time.Duration // 采样窗口的持续时间
	WinBucket    int           // 表示桶的数量
	CPUThreshold int64         // CPU的阈值
}

func (b *BBR) maxPASS() int64 {

	rawMaxPass := atomic.LoadInt64(&b.rawMaxPass)

	// 若有了rawMaxPass，并且 还没有开始新一次的采样,那么就会使用上一次的maxPass
	if rawMaxPass > 0 && b.passStat.Timespan() < 1 {
		return rawMaxPass
	}

	rawMaxPass = int64(b.passStat.Reduce(func(iterator metric.Iterator) float64 {
		var result = 1.0
		// 这里从
		for i := 1; iterator.Next() && i < b.conf.WinBucket; i++ {
			bucket := iterator.Bucket()

			count := 0.0

			for _, p := range bucket.Points {
				count += p
			}

			result = math.Max(result, count)
		}
		return result

	}))

	if rawMaxPass == 0 {
		rawMaxPass = 1
	}

	atomic.StoreInt64(&b.rawMaxPass, rawMaxPass)

	return rawMaxPass
}

func (b *BBR) minRT() int64 {

	rawMinRt := atomic.LoadInt64(&b.rawMinRt)

	if rawMinRt > 0 && b.rtStat.Timespan() < 1 {
		return rawMinRt
	}

	rawMinRt := int64(b.rtStat.Reduce(func(iterator metric.Iterator) float64 {

		var result = math.MaxFloat64

		for i := 1 ;

	}))



}

func (b *BBR) maxFlight() int64 {
	return 0
}

// 返回true 表示本次请求应该丢弃， 返回false表示本次请求不该丢弃
func (b *BBR) shouldDrop() bool {

	// 若当前cpu使用率还没有达到阈值
	if b.cpu() < b.conf.CPUThreshold {
		// 则去查看 上次丢弃请求的时间
		prevDrop, _ := b.prevDrop.Load().(time.Duration)
		// 若没有丢弃过请求 那么本次请求就不该丢弃
		if prevDrop == 0 {
			return false
		}
		// 若两次时间 相差 太短 那么就要 判断是否是瞬时流量太大了
		if time.Since(initTime)-prevDrop <= time.Second {
			inFlight := atomic.LoadInt64(&b.inFlight)
			return inFlight > 1 && inFlight > b.maxFlight()
		}
		b.prevDrop.Store(time.Duration(0))
		return false
	}
	// 若cpu的值已经超过阈值 那么就来判断当前处理的请求数量 是不是超过最大的处理数量
	inFlight := atomic.LoadInt64(&b.inFlight)
	drop := inFlight > 1 && inFlight > b.maxFlight()

	if drop {
		prevDrop, _ := b.prevDrop.Load().(time.Duration)

		if prevDrop != 0 {
			return drop
		}

		b.prevDrop.Store(time.Since(initTime))
	}
	return drop
}

func (b *BBR) Allow() (func(info DoneInfo), error) {

}

func newLimiter(conf *BBRConfig) Limiter {

	if conf == nil {
		conf = defaultConf
	}

	size := conf.WinBucket
	bucketDuration := conf.Window / time.Duration(conf.WinBucket)

	bbr := &BBR{}

	bbr.passStat = metric.NewRollingCounter(size, bucketDuration)
	bbr.rtStat = metric.NewRollingCounter(size, bucketDuration)
	bbr.conf = conf
	bbr.winBucketPerSec = int64(time.Second) / int64(conf.Window) / int64(conf.WinBucket)

	return bbr
}
