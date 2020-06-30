package oauth

import (
	"testing"
	"tomm/log"
)

func TestDB(t *testing.T) {
	secretInfo, err := getSecretInfo("qweqw")
	if err != nil {

		return
	}

	log.Info("secretInfo secret key is %s , appkey is %s", secretInfo.SecretKey, secretInfo.AppKey)

}

func TestCreateOAuthInfo(t *testing.T) {
	secretInfo, err := CreateOAuthInfo("")
	if err != nil {
		log.Error("Create Auth Info Fail Err is %s", err.Error())
		return
	}

	log.Info("SecretInfo AppKey is %s ,SecretKey is %s", secretInfo.AppKey, secretInfo.SecretKey)
}
