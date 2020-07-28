package service

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"tomm/api/api"
	"tomm/api/model"
	"tomm/core/server"
	"tomm/ecode"
	"tomm/log"
	"tomm/redis"
	"tomm/service/dao"
	"tomm/utils"
)

var (
	vcode = []string{"a", "b", "c", "d", "e", "f", "g", "h", "k", "l", "m", "n", "p", "q", "r",
		"t", "w", "y", "z", "2", "3", "4", "5", "7", "8", "A", "B", "C", "D", "E", "F", "G", "H",
		"K", "L", "M", "N", "P", "Q", "R", "T", "W", "Y", "Z"}
)

func RegisterAdmin(c *server.Context) {
	req := api.RegisterAdminReq{}
	err := c.Bind(&req)
	if err != nil {
		log.Warn("RegisterAdmin Bind Fail err is %s", err.Error())
		server.HttpCode(c, ecode.ParamFail)
		return
	}

	// 这里可以直接将 login_name 存放到redis中
	loginInfo := model.AdminInfos{
		LoginName: req.LoginName,
		LoginPwd:  req.LoginPwd,
		Name:      req.Name,
		Number:    req.Number,
	}

	loginInfo.LoginPwd = utils.Base64Encode([]byte(loginInfo.LoginPwd))
	err = dao.SaveAdminLogin(loginInfo)

	if err != nil {
		log.Error("RegisterAdmin SaveAdminLogin Fail err is %s", err.Error())
		server.HttpCode(c, ecode.SystemFail)
		return
	}
	res := api.RegisterAdminRes{
		Res: 1,
	}

	server.HttpData(c, res)
}

func GetVerificationCode(c *server.Context) {

	code := createCode(4)
	randomV, _ := utils.StrUUID()
	res := api.VCodeRes{
		Code:        code,
		RandomValue: randomV,
	}

	// 保存到redis中
	err := redis.Set(context.TODO(), fmt.Sprintf(dao.VCODE, randomV), code, dao.VCODE_EXP)

	if err != nil {
		log.Error("GetVerificationCode Set Redis Fail err is %s", err.Error())
		server.HttpCode(c, ecode.SystemFail)
		return
	}

	server.HttpData(c, res)
}

func AdminLogin(c *server.Context) {
	req := api.AdminLoginReq{}
	err := c.Bind(&req)

	if err != nil {
		log.Warn("AdminLogin Bind Fail err is %s", err.Error())
		server.HttpCode(c, ecode.ParamFail)
		return
	}

	// 检查 VCode
	vCode := ""
	key := fmt.Sprintf(dao.VCODE, req.Random)
	err = redis.Get(context.TODO(), key, &vCode)
	if err != nil && ecode.NotValue.EqualErr(err) {
		log.Warn("AdminLogin Get VCode Fail key is %s err is %s", key, err.Error())
		server.HttpCode(c, ecode.LoginFail)
		return
	}

	if vCode != req.VCode {
		log.Warn("AdminLogin VCode Wrong vCode is %s , req VCode is %s", vCode, req.VCode)
		server.HttpCode(c, ecode.VCodeFail)
		return
	}

	// 检查用户名
	exist, err := dao.CheckAdminLoginName(req.LoginName)

	if err != nil {
		// VCodeFail
		log.Warn("AdminLogin CheckAdminLoginName Fail LoginName is %s , err  is %s", req.LoginName, err.Error())
		server.HttpCode(c, ecode.SystemFail)
		return
	}

	if !exist {
		log.Debug("AdminLogin LoginName Not Exist LoginName is %s", req.LoginName)
		server.HttpCode(c, ecode.LoginFail)
		return
	}

	// 获取密码
	pwd := utils.Base64Encode([]byte(req.LoginPwd))
	adminInfo, err := dao.GetAdminInfoByLoginName(req.LoginName)

	if err != nil {
		// VCodeFail
		log.Warn("AdminLogin GetPwd Fail LoginName is %s , err  is %s", req.LoginName, err.Error())
		server.HttpCode(c, ecode.LoginFail)
		return
	}
	// 检查密码
	if pwd != adminInfo.LoginPwd {
		log.Debug("AdminLogin Pwd Wrong loginName is %s ,pwd is %s", req.LoginName, pwd)
		server.HttpCode(c, ecode.LoginFail)
		return
	}

	redis.Del(context.TODO(), key)

	server.HttpData(c, adminInfo)

}

func AdminUpdatePwdSafe(c *server.Context) {

	req := api.AdminUpdatePwdSafeReq{}

	err := c.Bind(&req)

	if err != nil {
		log.Warn("AdminUpdatePwdSafe Bind Fail err is %s", err.Error())
		server.HttpCode(c, ecode.ParamFail)
		return
	}

	if req.NewPwd == req.OldPwd {
		log.Debug("AdminUpdatePwdSafe NewPwd==OldPwd")
		server.HttpCode(c, ecode.PwdEqualFail)
	}

	newPwd := utils.Base64Encode([]byte(req.NewPwd))
	oldPwd := utils.Base64Encode([]byte(req.OldPwd))

	err = dao.UpdatePwdByLoginNameSafe(req.LoginName, newPwd, oldPwd)

	if err != nil {
		log.Warn("AdminUpdatePwdSafe UpdatePwdByLoginNameSafe Fail err is %s", err.Error())
		server.HttpCode(c, ecode.EditFail)
		return
	}

	server.HttpCode(c, nil)
}

func AdminUpdatePwd(c *server.Context) {

	req := api.AdminUpdatePwdReq{}

	err := c.Bind(&req)

	if err != nil {
		log.Warn("AdminUpdatePwd Bind Fail err is %s", err.Error())
		server.HttpCode(c, ecode.ParamFail)
		return
	}
	newPwd := utils.Base64Encode([]byte(req.NewPwd))
	err = dao.UpdatePwdByLoginName(req.LoginName, newPwd)

	if err != nil {
		log.Warn("AdminUpdatePwdSafe UpdatePwdByLoginNameSafe Fail err is %s", err.Error())
		server.HttpCode(c, ecode.EditFail)
		return
	}

	server.HttpCode(c, nil)
}

//
//func DeleteAdmin(c *server.Context) {
//
//
//
//}

func CheckAdminName(c *server.Context) {
	// 直接从redis里面取就好了
	req := api.CheckAdminNameReq{}

	err := c.Bind(&req)

	if err != nil {
		log.Warn("CheckAdminName Bind Fail err is %s", err.Error())
		server.HttpCode(c, ecode.ParamFail)
		return
	}

	exist, err := dao.CheckAdminLoginName(req.LoginName)

	if err != nil {
		log.Warn("RegisterAdmin CheckAdminLoginName Fail err is %s", err.Error())
		server.HttpCode(c, ecode.ParamFail)
		return
	}
	res := &api.CheckAdminNameRes{}
	if exist {
		// 创建
		res.Res = 2
	} else {
		res.Res = 1
	}

	server.HttpData(c, res)

}

func createCode(num int) string {

	build := strings.Builder{}
	rand.Seed(time.Now().Unix())
	for i := 0; i < num; i++ {
		build.WriteString(vcode[rand.Intn(len(vcode))])
	}

	return build.String()
}
