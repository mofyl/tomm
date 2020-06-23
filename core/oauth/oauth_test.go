package oauth

import (
	"go.uber.org/zap"
	"testing"
	"tomm/log"
)

func TestDB(t *testing.T) {
	secretInfo, err := getSecretInfo("qweqw")
	if err != nil {
		log.Msg(log.INFO, err.Error())
		return
	}

	log.Info("secretInfo", zap.String("secretKey", secretInfo.SecretKey), zap.String("appKey", secretInfo.AppKey))

}
