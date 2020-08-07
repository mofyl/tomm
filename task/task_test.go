package task

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func TestAtomic(t *testing.T) {

	var tmp int32
	tmp = 1
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
	for i := 0; i < 10000; i++ {
		ctx, err := NewTaskContext(c, 2, true, func(ctx *TaskContext) bool {
			//fmt.Println(111)
			//time.Sleep(3 * time.Second)
			return true
		}, func(ctx *TaskContext) bool {
			//fmt.Println(2222)
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

	wg.Add(1)
	go func() {
		fmt.Println("Wait Res")
		for {
			select {
			case v, ok := <-c:
				if !ok {
					for v := range c {
						fmt.Println(v.TaskID)
					}
					wg.Done()
					fmt.Println("Finish")
					return
				}
				fmt.Println(v.TaskID)
			}

		}

	}()

	wg.Add(1)
	go func() {
		tm.Close()
		wg.Done()
		close(c)
	}()
	wg.Wait()

}

func BenchmarkTask(b *testing.B) {

	//c := make(chan *TaskContext)
	//wg := &sync.WaitGroup{}
	initNoConf()

	b.ResetTimer()

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

	c := make(chan int, 3)

	c <- 2
	c <- 3
	close(c)

	v, ok := <-c

	fmt.Printf("%d , %v\n", v, ok)

	for v := range c {

		fmt.Printf("%d\n", v)
	}

}
