package service

import (
	"context"
	"io/ioutil"
	"net/http"
	"tomm/api/service"
	"tomm/ecode"
	"tomm/log"
	"tomm/redis"
	"tomm/service/dao"
	"tomm/task"
	"tomm/utils"
)

func getTokenJob1(ctx *task.TaskContext) bool {
	// *service.PlatformInfo, *service.ReqDataInfo,
	var platformInfo *service.PlatformInfo
	var reqDataInfo *service.TokenDataInfo
	if plat, ok := ctx.Get("secretInfo"); ok {
		platformInfo, _ = plat.(*service.PlatformInfo)
	}

	if reqData, ok := ctx.Get("reqDataInfo"); ok {
		reqDataInfo, _ = reqData.(*service.TokenDataInfo)
	}

	mmUserInfo, errMsg := GetUserInfo()

	if errMsg != nil {
		ctx.Err = errMsg
		return true
	}

	token, expTime, err := dao.GetToken(platformInfo.AppKey)
	if err != nil {
		log.Error("Get Token Fail err is %s", err.Error())
		ctx.Err = errMsg
		return true
	}
	// 关联token 和 userID
	redis.Set(context.TODO(), token, mmUserInfo, dao.RESOURCE_TOKEN_EXP)

	log.Debug("Return Token is %s", token)
	// 构造返回值
	// 返回值包括 token + expTime + extendInfo
	tokenInfo := service.TokenInfo{
		Token:      token,
		ExpTime:    expTime,
		ExtendInfo: reqDataInfo.ExtendInfo,
	}
	tokenB, _ := utils.Json.Marshal(tokenInfo)

	resBase64Str, err := utils.AESCBCBase64Encode(platformInfo.SecretKey, tokenB)
	if err != nil {
		log.Error("AESCBCBase64Encode Fail Err is %s", err.Error())
		ctx.Err = ecode.NewErr(err)
		return true
	}
	res := service.GetTokenRes{}
	//res.TokenInfo = resBase64Str
	res.Token = resBase64Str
	ctx.Err = nil
	ctx.Set("res", res)
	return true
	//httpData(c, res)
}

func getTokenJob2(ctx *task.TaskContext) bool {
	var platformInfo *service.PlatformInfo
	var tokenRes *service.GetTokenRes
	if plat, ok := ctx.Get("secretInfo"); ok {
		platformInfo, _ = plat.(*service.PlatformInfo)
	}
	//
	//if reqData, ok := ctx.Get("reqDataInfo"); ok {
	//	reqDataInfo, _ = reqData.(*service.ReqDataInfo)
	//}

	if resInterface, ok := ctx.Get("res"); ok {
		tokenRes = resInterface.(*service.GetTokenRes)
	}

	var rsp *http.Response
	var err error
	if ctx.Err != nil {
		rsp, err = backCode(platformInfo.SignUrl, ctx.Err)
	} else {
		rsp, err = backData(platformInfo.SignUrl, tokenRes)
	}

	if err != nil {
		// 重发机制
	}

	rspB, _ := ioutil.ReadAll(rsp.Body)
	if string(rspB) != "1" {
		// 重复 发送
	}

	return false

}

//
//func NewJobUserInfo(err ecode.ErrMsgs, info *job.JobUserInfo) *task.TaskOut {
//
//	b, _ := info.Marshal()
//	return task.NewTaskOut(task.GetUserJob, err, b)
//}
//
//func (s *Ser) job() {
//	for {
//		select {
//		case jobRes, ok := <-s.jobNotify:
//			if !ok {
//				s.wg.Done()
//				return
//			}
//
//			switch jobRes.Type {
//			case task.GetUserJob:
//				s.getUserJob(jobRes)
//			}
//
//		}
//	}
//}
//
//func (s *Ser) getUserJob(res *task.TaskOut) {
//
//	//job := job.JobUserInfo{}
//	//
//	//err := job.Unmarshal(res.Data)
//	//
//	//if err != nil {
//	//	log.Error("Do GetUserInfo PoolJob : Unmarshal Fail ")
//	//	return
//	//}
//
//	//s.p.DoJob(&job.Job{
//	//	ID:        1111,
//	//	ResNotify: s.jobNotify,
//	//	Do: func() *job.JobRes {
//	//
//	//		notifyRes := &rending.Json{}
//	//
//	//		notifyRes.Code = res.Err.Code()
//	//		notifyRes.Msg = res.Err.Error()
//	//
//	//		if job.Base64Str != "" {
//	//			notifyRes.Data = job.Base64Str
//	//		}
//	//
//	//		rsp, err := HttpJsonPost(job.CallBack, notifyRes)
//	//		defer rsp.Body.Close()
//	//		if err != nil {
//	//			log.Error("UseInfo Notify Third Part Fail CallBack is %s , err is %s", job.CallBack, err.Error())
//	//			return job.NewJobRes(job.JobFail, ecode.NewErr(err), nil)
//	//		}
//	//
//	//		if rsp.StatusCode != 200 {
//	//			log.Error("UseInfo Notify Third Part Fail CallBack is %s,return statusCode is %d", job.CallBack, rsp.StatusCode)
//	//			return job.NewJobRes(job.JobFail, ecode.NewErr(errors.New(fmt.Sprintf("Status Code is %d ", rsp.StatusCode))), nil)
//	//		}

//	//
//	//	},
//	//})
//
//}
