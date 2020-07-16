package service

import (
	"tomm/config"
	"tomm/core/server"
)

const (
	DATALEN = 4

	TIMELEN          = 4
	CHANNEL_INFO_LEN = 4
	MAX_DATA         = 512

	MAX_TTL = 5 * 60 // 若数据包: nowTime - sendTime > 10min 则不处理
)

var (
	defaultConf *ServiceConf
)

func init() {
	defaultConf = &ServiceConf{}

	if err := config.Decode(config.CONFIG_FILE_NAME, "server", defaultConf); err != nil {
		panic("Service Load Config Fail " + err.Error())
	}

	if err := config.Decode(config.CONFIG_FILE_NAME, "mmServer", &mmSerUrl); err != nil {
		panic("Service Load mm Server Url Fail " + err.Error())
	}
}

type ServiceConf struct {
	NotifyChan int64 `yaml:"notifyChan"`
}

type Ser struct {
	e    *server.Engine
	conf *ServiceConf
	//jobNotify chan *task.TaskContext
	//wg *sync.WaitGroup
	//p         *task.TaskManager

	tokenGroup    server.IRouter
	platformGroup server.IRouter
}

func NewService() *Ser {
	s := &Ser{
		//jobNotify: make(chan *task.TaskContext, defaultConf.NotifyChan),
		//wg: &sync.WaitGroup{},
	}
	e := server.NewEngine(nil)
	s.e = e
	s.conf = defaultConf
	s.platformGroup = s.e.NewGroup("/platform")
	s.tokenGroup = s.e.NewGroup("/token")
	s.registerRouter()

	return s
}

func (s *Ser) registerRouter() {
	s.tokenGroup.GET("/getToken", GetResourceToken)
	s.tokenGroup.GET("/verifyToken", VerifyToken)
	s.tokenGroup.GET("/getUserInfo", GetUserInfo)
	s.tokenGroup.GET("/getCode", GetCode)
	s.tokenGroup.GET("/checkCode", CheckCode)
	// checkCode Appkey+Code
	//
	s.platformGroup.POST("/register", RegisterPlatform)
	s.platformGroup.GET("/checkPlatformName", CheckPlatformName)

	//s.userGroup.GET("/getCode", s.getCode)

}

func (s *Ser) Close() {
	s.e.Close()

	//task.Close()
	//close(s.jobNotify)

	//s.wg.Wait()
}

func (s *Ser) Start() {
	s.e.RunServer()
}
