package service

import (
	"testing"
	"tomm/log"
	"tomm/service/dao"
)

func TestDB(t *testing.T) {
	secretInfo, err := dao.GetPlatformInfo("qweqw")
	if err != nil {

		return
	}

	log.Info("secretInfo secret key is %s , appkey is %s", secretInfo.SecretKey, secretInfo.AppKey)

}

//
//func TestCreateOAuthInfo(t *testing.T) {
//	secretInfo, err := CreateOAuthInfo("")
//	if err != nil {
//		log.Error("Create Auth Info Fail Err is %s", err.Error())
//		return
//	}
//
//	log.Info("PlatformInfo AppKey is %s ,SecretKey is %s", secretInfo.AppKey, secretInfo.SecretKey)
//}
