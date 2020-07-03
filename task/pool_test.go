package task

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	//res := make(chan []byte, 100)
	for i := 0; i < 100; i++ {
		id := int64(i)
		job := &PoolJob{
			ID:        id,
			ResNotify: nil,
			Do: func() []byte {
				//time.Sleep(3 * time.Second)
				ids := strconv.FormatInt(id, 10)
				return []byte(ids)
			},
		}
		DoJob(job)
	}
	fmt.Println("Send Finish")
	Close()
}

func BenchmarkDoJob(b *testing.B) {
	//b.ResetTimer()
	time.Sleep(2 * time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := int64(i)
		job := &PoolJob{
			ID:        id,
			ResNotify: nil,
			Do: func() []byte {
				//time.Sleep(3 * time.Second)
				ids := strconv.FormatInt(id, 10)
				return []byte(ids)
			},
		}
		DoJob(job)
	}
	fmt.Println("Send Finish")
	Close()
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
