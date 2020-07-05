package service

import (
	"context"
	"sync"
	"time"
	"tomm/api/job"
	"tomm/config"
	"tomm/core/server"
	"tomm/ecode"
	"tomm/log"
	"tomm/redis"
	"tomm/service/oauth"
	"tomm/task"
	"tomm/utils"
)

var (
	defaultConf *ServiceConf
)

func init() {
	defaultConf = &ServiceConf{}

	if err := config.Decode(config.CONFIG_FILE_NAME, "server", defaultConf); err != nil {
		panic("Service Load Config Fail " + err.Error())
	}
}

type ServiceConf struct {
	NotifyChan int64 `yaml:"notifyChan"`
}

type Ser struct {
	e         *server.Engine
	conf      *ServiceConf
	jobNotify chan *task.TaskOut
	wg        *sync.WaitGroup
	p         *task.Pool
}

func NewService() *Ser {
	s := &Ser{
		jobNotify: make(chan *task.TaskOut, defaultConf.NotifyChan),
		wg:        &sync.WaitGroup{},
	}
	e := server.NewEngine(nil)
	s.e = e
	s.conf = defaultConf
	s.p = task.NewPool(nil, s.wg)

	s.registerRouter()
	return s
}

func (s *Ser) registerRouter() {
	s.e.GET("/getToken", s.getResourceToken)
	s.e.GET("/verifyToken", s.verifyToken)
}

func (s *Ser) Close() {
	s.e.Close()
	s.p.Close()
	cli.CloseIdleConnections()
	close(s.jobNotify)

	s.wg.Wait()
}

func (s *Ser) Start() {
	s.e.RunServer()

	s.wg.Add(1)

	go s.job()
}

func (s *Ser) verifyToken(c *server.Context) {
	req := VerifyTokenReq{}
	err := c.Bind(&req)
	if err != nil {
		log.Warn("VerifyToken Bind Err is %s", err.Error())
	}

	// 查看token是否存在
	token, expTime, err := oauth.GetToken(req.AppKey)

	if err != nil {
		log.Error("Verify Token Fail Err is %s", err.Error())
		httpCode(c, ecode.ParamFail)
	}

	if token != token {
		// 校验失败
		httpCode(c, ecode.VerifyFail)
	}
	res := VerifyTokenRes{}
	res.ExpTime = expTime
	httpData(c, res)
}

func (s *Ser) getResourceToken(c *server.Context) {

	secretInfo, reqDataInfo, eCode := checkGetTokenReq(c)
	if eCode != nil {
		httpCode(c, eCode)
	}

	httpCode(c, ecode.OK)

	// 解密完成 第三方等待回调
	// 到资源服务器请求 查看是否授权
	s.p.DoJob(&task.PoolJob{
		ID:        111,
		ResNotify: s.jobNotify,
		Do: func() *task.TaskOut {
			//
			jobInfo := &job.JobUserInfo{
				CallBack: reqDataInfo.BackUrl,
			}
			mmUserInfo, errMsg := GetUserInfo()
			if errMsg != nil {
				return NewJobUserInfo(errMsg, jobInfo)
			}

			token, expTime, err := oauth.GetToken(secretInfo.AppKey)
			if err != nil {
				log.Error("Get Token Fail err is %s", err.Error())
				return NewJobUserInfo(errMsg, jobInfo)
			}
			// 关联token 和 userID
			redis.Set(context.TODO(), token, mmUserInfo, oauth.RESOURCE_TOKEN_EXP)

			log.Debug("Return Token is %s", token)
			// 构造返回值
			// 返回值包括 token + expTime + extendInfo
			tokenInfo := TokenInfo{
				Token:      token,
				ExpTime:    expTime,
				ExtendInfo: reqDataInfo.ExtendInfo,
			}
			tokenB, _ := utils.Json.Marshal(tokenInfo)

			resBase64Str, err := utils.AESCBCBase64Encode(secretInfo.SecretKey, tokenB)
			if err != nil {
				log.Error("AESCBCBase64Encode Fail Err is %s", err.Error())
				return NewJobUserInfo(errMsg, jobInfo)
			}
			//res := GetTokenRes{}
			//res.TokenInfo = resBase64Str
			jobInfo.Base64Str = resBase64Str
			return NewJobUserInfo(ecode.OK, jobInfo)
			//httpData(c, res)
		},
	})
	return
}

func checkGetTokenReq(c *server.Context) (*oauth.SecretInfo, *ReqDataInfo, ecode.ErrMsgs) {
	req := &GetTokenReq{}
	err := c.Bind(req)

	if err != nil {
		log.Warn("GetToken Bind Err is %s ", err.Error())
		return nil, nil, ecode.ParamFail
	}

	if req.AppKey == "" || req.Data == "" {
		return nil, nil, ecode.ParamFail
	}

	// 获取该appKey
	secretInfo, err := oauth.GetOAuthInfo(req.AppKey)
	if err != nil || secretInfo == nil {
		if err != nil {
			log.Error("GetToken Fail AppKey is %s , Err is %s", req.AppKey, err.Error())
		}
		return nil, nil, ecode.SecretKeyFail
	}

	reqDataInfo, eCode := GetDataInfo(secretInfo.SecretKey, req.Data)
	if eCode != nil {
		log.Error("Get Data Info Fail ")
		return nil, nil, eCode
	}

	if reqDataInfo.ChannelInfo == "" ||
		reqDataInfo.SendTime == 0 ||
		!utils.CheckUrl(reqDataInfo.BackUrl) {
		return nil, nil, ecode.ParamFail
	}
	// 超过10分钟就不处理了
	if time.Now().Unix()-int64(reqDataInfo.SendTime) > MAX_TTL {
		log.Error("Package Timeout ")
		return nil, nil, ecode.PackageTimeout
	}

	return secretInfo, reqDataInfo, nil
}

func getResourceUrl(resourceUrl string) string {
	return resourceUrl + ""
}

//
//func buildRes(c *server.Context, res *GetTokenRes, secretKey string, code int) {
//	var b []byte
//	if res != nil {
//		b, _ := json.Marshal(res)
//		origData, err := utils.AESCBCBase64Encode(secretKey, b)
//		if err != nil {
//			code = 500
//			b = nil
//		} else {
//			b = []byte(origData)
//		}
//	} else {
//		code = 400
//		b = nil
//	}
//
//	c.Json(code, b)
//}

func GetDataInfo(secretKey string, data string) (*ReqDataInfo, ecode.ErrMsgs) {

	// 使用 secretKey 进行 AES解密
	origData, err := utils.AESCBCBase64Decode(secretKey, data)

	if err != nil {
		return nil, ecode.DeCodeFail
	}

	origLen := len(origData)

	if origLen < DATALEN+TIMELEN+CHANNEL_INFO_LEN || origLen > MAX_DATA {
		return nil, ecode.ParamFail
	}
	//
	//dataLen := binary.BigEndian.Uint32(origData[:DATALEN])
	//sendTime := binary.BigEndian.Uint64(origData[DATALEN : DATALEN+TIMELEN])
	//channelInfo := origData[DATALEN+TIMELEN : DATALEN+TIMELEN+CHANNEL_INFO_LEN]
	//extendInfo := origData[DATALEN+TIMELEN+CHANNEL_INFO_LEN:]
	//reqInfo := &ReqDataInfo{
	//	DataLen:     int32(dataLen),
	//	SendTime:    int64(sendTime),
	//	ChannelInfo: string(channelInfo),
	//	ExtendInfo:  extendInfo,
	//}

	reqInfo := &ReqDataInfo{}
	err = utils.Json.Unmarshal(origData, reqInfo)

	if err != nil {
		return nil, ecode.ParamFail
	}
	//if int(reqInfo.DataLen) != origLen {
	//	return nil, DECODE_FAIL, DECODE_FAIL_MSG
	//}
	return reqInfo, nil
}
