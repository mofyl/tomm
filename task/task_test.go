package task

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestAtomic(t *testing.T) {

	var tmp int32 = 1
	if atomic.CompareAndSwapInt32(&tmp, 1, 3) {
		fmt.Println(11)
		return
	}

	fmt.Println(2)
}

func TestTask(t *testing.T) {

	c := make(chan *TaskContext, 10)
	wg := &sync.WaitGroup{}
	initNoConf()
	//

	wg.Add(1)
	go func() {
		fmt.Println("Wait Res")
		for v := range c {
			fmt.Println(v.TaskID)
		}

		wg.Done()
	}()

	wg.Add(1)
	go func() {
		time.Sleep(1 * time.Second)
		tm.Close()
		wg.Done()
		close(c)
	}()

	for i := 0; i < 1000; i++ {
		ctx, err := NewTaskContext(c, 2, true, func(ctx *TaskContext) bool {
			fmt.Println(111)
			//time.Sleep(3 * time.Second)
			return true
		}, func(ctx *TaskContext) bool {
			fmt.Println(2222)
			//time.Sleep(3 * time.Second)
			return true
		})

		if err != nil {
			fmt.Println("err ", err.Error())
			return
		}

		res := ctx.Start()

		fmt.Printf("Ctx Start Res is %v , index is %d\n", res, i)
	}

	wg.Wait()
}

func TestWaitRes(t *testing.T) {

	num := 9
	c := make(chan int, num)

	go func() {

		for i := 0; i < num; i++ {
			c <- i
		}

	}()

	for {
		if len(c) == num {
			fmt.Println("finish")
			break
		}
	}

}

func BenchmarkTask(b *testing.B) {

	//c := make(chan *TaskContext)
	//wg := &sync.WaitGroup{}
	initNoConf()

	b.ResetTimer()

	go func() {
		time.Sleep(2 * time.Second)
		tm.Close()
	}()

	for i := 0; i < b.N; i++ {
		ctx, err := NewTaskContext(nil, 2, true, func(ctx *TaskContext) bool {
			//time.Sleep(1 * time.Second)
			fmt.Println(111)
			return true
		}, func(ctx *TaskContext) bool {
			//time.Sleep(1 * time.Second)
			fmt.Println(222)
			return true
		})

		if err != nil {
			fmt.Println("err ", err.Error())
			return
		}

		ctx.Start()
	}

	//wg.Wait()

}

func TestChannel(t *testing.T) {

	c := make(chan int)

	go func() {
		v, ok := <-c
		fmt.Println("111 ", v, ok)
	}()

	go func() {
		v, ok := <-c
		fmt.Println("222  ", v, ok)
	}()
	time.Sleep(2 * time.Second)
	c <- 1

	close(c)
	select {}

}
