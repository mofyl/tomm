package metric

import (
	"sync"
	"time"
)

type RollingPolicy struct {
	mu     sync.RWMutex
	size   int
	window *Window
	offset int

	// 一个 bucket的持续时间
	bucketDuration time.Duration
	// 最后一次添加的时间
	lastAppendTime time.Time
}

func NewRollingPolicy(window *Window, bucketDuration time.Duration) *RollingPolicy {
	return &RollingPolicy{
		mu:             sync.RWMutex{},
		size:           window.Size(),
		window:         window,
		offset:         0,
		bucketDuration: bucketDuration,
		lastAppendTime: time.Now(),
	}
}

func (r *RollingPolicy) timespan() int {
	/*
		这里使用 上次添加的时间到现在的 这一个段时间 / 桶的持续时间 = 该经历多少个桶
		后面在使用offset + 就可以算出 当前时间应该在哪个桶
	*/
	v := int(time.Since(r.lastAppendTime) / r.bucketDuration)

	if v > -1 {
		return v
	}
	return r.size
}

func (r *RollingPolicy) add(f func(offset int, val float64), val float64) {
	// 修改 lastAppendTime
	// 获取 经历了多少个桶
	r.mu.Lock()
	timespan := r.timespan()
	// 若这段时间内经历过了桶
	if timespan > 0 {
		r.lastAppendTime = r.lastAppendTime.Add(time.Duration(timespan) * r.bucketDuration)

		offset := r.offset
		s := offset + 1

		if timespan > r.size {
			timespan = r.size
		}

		e, e1 := s+timespan, 0

		if e > r.size {
			e1 = e - r.size
			e = r.size
		}

		// 清空后面的
		for i := s; i < e; i++ {
			r.window.ReSetBucket(i)
			offset = i
		}

		for i := 0; i < e1; i++ {
			r.window.ReSetBucket(i)
			offset = i
		}

		r.offset = offset

	}

	f(r.offset, val)
	r.mu.Unlock()
}

func (r *RollingPolicy) Append(val float64) {
	r.add(r.window.Append, val)
}

func (r *RollingPolicy) Add(val float64) {
	r.add(r.window.Add, val)
}

func (r *RollingPolicy) Reduce(f func(Iterator) float64) float64 {
	r.mu.RLock()
	timespan := r.timespan()
	var val float64
	// 这里 表示 现在的时间距离上次采点的时间 大于 bucketDuration  不然 count为0
	if count := r.size - timespan; count > 0 {
		offset := timespan + r.offset + 1
		// 这里判断 当前位置是否已经超出了size， 若超出了size 就要返回头部 从头部再次计算
		if offset > r.size {
			offset = offset - r.size
		}
		val = f(r.window.Iterator(offset, count))
	}
	r.mu.RUnlock()
	return val

}
