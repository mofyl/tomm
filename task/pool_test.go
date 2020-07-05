package task

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var (
	test_wg *sync.WaitGroup
)

func newPool() *Pool{
	test_wg = &sync.WaitGroup{}
	p := NewPool(nil , test_wg)
	return p
}

func TestPool(t *testing.T) {
	//res := make(chan []byte, 100)
	p := newPool()
	for i := 0; i < 100; i++ {
		id := int64(i)
		job := &Job{
			ID:        id,
			ResNotify: nil,
			Do: func() *TaskContext {
				//time.Sleep(3 * time.Second)
				//ids := strconv.FormatInt(id, 10)
				return nil
			},
		}
		p.DoJob(job)
	}
	fmt.Println("Send Finish")
	p.Close()
	test_wg.Wait()
}

func BenchmarkDoJob(b *testing.B) {
	//b.ResetTimer()
	p := newPool()
	time.Sleep(2 * time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := int64(i)
		job := &Job{
			ID:        id,
			ResNotify: nil,
			Do: func() *TaskContext {
				//time.Sleep(3 * time.Second)
				//ids := strconv.FormatInt(id, 10)
				return nil
			},
		}
		p.DoJob(job)
	}
	fmt.Println("Send Finish")
	p.Close()
	test_wg.Wait()
}

func TestSelect(t *testing.T) {

	c := make(chan int)

	close(c)

	select {
	case c <- 1:
	default:
		fmt.Println(1111)
	}

	fmt.Println(2222)
}
