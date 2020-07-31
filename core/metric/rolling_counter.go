package metric

import (
	"fmt"
	"time"
)

type RollingCounter struct {
	policy *RollingPolicy
}

func NewRollingCounter(size int, bucketDuration time.Duration) *RollingCounter {

	window := NewWindow(size)
	policy := NewRollingPolicy(window, bucketDuration)

	return &RollingCounter{
		policy: policy,
	}

}

func (r *RollingCounter) Add(val int64) {
	if val < 0 {
		panic(fmt.Errorf("stat/metric: cannot decrease in value. val: %d", val))
	}
	r.policy.Add(float64(val))
}

func (r *RollingCounter) Reduce(f func(Iterator) float64) float64 {
	return r.policy.Reduce(f)
}

func (r *RollingCounter) Timespan() int {
	return r.policy.timespan()
}
