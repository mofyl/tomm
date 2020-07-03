package service

import (
	"tomm/ecode"
	"tomm/log"
	"tomm/service/api"
	"tomm/task"
)

func NewJobUserInfo(err ecode.ErrMsgs, info *api.JobUserInfo) *task.TaskOut {

	b, _ := info.Marshal()
	return task.NewTaskOut(task.GetUserJob, err, b)
}

func (s *Ser) job() {
	for {
		select {
		case jobRes, ok := <-s.jobNotify:
			if !ok {
				s.wg.Done()
				return
			}

			switch jobRes.Type {
			case task.GetUserJob:
				s.getUserJob(jobRes)
			}

		}
	}
}

func (s *Ser) getUserJob(res *task.TaskOut) {

	job := api.JobUserInfo{}

	err := job.Unmarshal(res.Data)

	if err != nil {
		log.Error("Do GetUserInfo PoolJob : Unmarshal Fail ")
		return
	}

	//s.p.DoJob(&job.Job{
	//	ID:        1111,
	//	ResNotify: s.jobNotify,
	//	Do: func() *job.JobRes {
	//
	//		notifyRes := &rending.Json{}
	//
	//		notifyRes.Code = res.Err.Code()
	//		notifyRes.Msg = res.Err.Error()
	//
	//		if job.Base64Str != "" {
	//			notifyRes.Data = job.Base64Str
	//		}
	//
	//		rsp, err := HttpJsonPost(job.CallBack, notifyRes)
	//		defer rsp.Body.Close()
	//		if err != nil {
	//			log.Error("UseInfo Notify Third Part Fail CallBack is %s , err is %s", job.CallBack, err.Error())
	//			return job.NewJobRes(job.JobFail, ecode.NewErr(err), nil)
	//		}
	//
	//		if rsp.StatusCode != 200 {
	//			log.Error("UseInfo Notify Third Part Fail CallBack is %s,return statusCode is %d", job.CallBack, rsp.StatusCode)
	//			return job.NewJobRes(job.JobFail, ecode.NewErr(errors.New(fmt.Sprintf("Status Code is %d ", rsp.StatusCode))), nil)
	//		}
	//
	//		rspB, _ := ioutil.ReadAll(rsp.Body)
	//		if string(rspB) == "1" {
	//			// 重复 发送
	//		}
	//
	//	},
	//})

}
