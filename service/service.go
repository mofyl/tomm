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

	token    server.IRouter
	platform server.IRouter
	admin    server.IRouter // 保存第三方平台管理员接口
	auth     server.IRouter // 保存 mm用户对第三方平台的 权限
}

func NewService() *Ser {
	s := &Ser{
		//jobNotify: make(chan *task.TaskContext, defaultConf.NotifyChan),
		//wg: &sync.WaitGroup{},
	}
	e := server.NewEngine(nil)
	s.e = e
	s.conf = defaultConf
	s.platform = s.e.NewGroup("/platform")
	s.token = s.e.NewGroup("/token")
	s.admin = s.e.NewGroup("/admin")
	s.auth = s.e.NewGroup("/auth")
	s.registerRouter()

	return s
}

func (s *Ser) registerRouter() {

	// 建立session
	// 小写 表示不对外公开的
	//s.e.GET("/startSession")

	s.token.GET("/token", GetResourceToken)
	s.token.GET("/verifyToken", VerifyToken)
	s.token.GET("/checkCode", CheckCode)

	s.token.GET("/UserInfo", GetUserInfo_V2)
	s.token.GET("/Code", GetCode)
	// checkCode Appkey+Code
	//
	s.platform.POST("/Register", RegisterPlatform)
	s.platform.GET("/CheckPlatformName", CheckPlatformName)
	s.platform.GET("/PlatformInfos", GetPlatformInfos)

	// 管理平台注册用户
	s.admin.POST("/Register", RegisterAdmin)
	// 管理平台用户获取验证码
	s.admin.GET("/VCode", GetVerificationCode)
	// 管理平台用户检查用户名是否存在
	s.admin.GET("/CheckLoginName", CheckAdminName)
	// 管理平台用户登录
	s.admin.POST("/Login", AdminLogin)
	//s.userGroup.GET("/getCode", s.getCode)
	// 设置权限组
	s.admin.POST("/PlatformRole", AddPlatformRole)

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
