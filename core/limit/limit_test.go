package limit

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
	"tomm/api/job"
	"tomm/task"
)

func TestGet(t *testing.T) {

	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.TODO())

	wg.Add(1)

	go put(30, ctx)

	// 模拟用户请求 一直过来拿桶

	ticker := time.NewTicker(1 * time.Minute)

	for {
		select {
		case <-ticker.C:
			ticker.Stop()
			task.Close()
			cancel()
			return
		default:
			_, c := task.NewTaskContext(nil, job.JobApi_JobUserInfo, 1,
				func(in *task.TaskContext) bool {
					Get()
					return false
				})
			c.Start()

		}
	}

}

func TestAllow(t *testing.T) {
	//
	window := 10 * time.Second

	winBucket := 100

	bucketDuration := window / time.Duration(winBucket)

	fmt.Println(bucketDuration.Milliseconds())

	tNow := time.Now()

	time.Sleep(1 * time.Second)
	vSince := time.Since(tNow)
	v := vSince / bucketDuration

	fmt.Println(vSince.Milliseconds())

	fmt.Println(int(v))

}
