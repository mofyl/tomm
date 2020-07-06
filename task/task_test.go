package task

import (
	"fmt"
	_ "net/http/pprof"
	"testing"
	"tomm/log"
)

//
func TestTask(t *testing.T) {

	notify := make(chan TaskContext)
	tm := NewTaskManager()

	md := NewTaskMD("11111",
		func(in *TaskContext) bool {
			fmt.Println("11111")
			in.Set("qwe", "qweqwewqe")
			return true
		},
		func(in *TaskContext) bool {
			fmt.Println("22222")
			fmt.Println(in.Get("qwe"))
			return true
		})

	ctx := tm.RegisterTaskMD(md, notify)
	ctx.Start()
	<-notify
	tm.Close()
	log.Info("TaskID is %s", ctx.TaskID)
}

func BenchmarkTaskManager(b *testing.B) {
	//
	//go func() {
	//
	//	if err := http.ListenAndServe(":10000", nil); err != nil {
	//		panic("ListenAndServer pprof Err " + err.Error())
	//	}
	//
	//}()

	notify := make(chan TaskContext)
	tm := NewTaskManager()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		md := TaskMD{
			TaskName:  fmt.Sprintf("1111_%d", i),
			TaskStage: 0,
			TaskHandlers: []TaskHandler{
				func(in *TaskContext) bool {
					fmt.Println("qweqwe")
					return false
				},
			},
		}

		ctx := tm.RegisterTaskMD(&md, notify)
		ctx.Start()
	}

	<-notify
	tm.Close()
	//select {}
}
