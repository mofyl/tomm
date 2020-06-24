package oauth

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
	"tomm/redis"
	"tomm/sqldb"
	"tomm/utils"
)

type SecretInfo struct {
	ChannelInfo string
	SecretKey   string
	AppKey      string
}

func SaveSecretInfo(info *SecretInfo) error {
	// save DB
	err := sqldb.GetDB(sqldb.MYSQL).Exec(context.TODO(), "insert into channel_infos(`channel_info` , `app_key` , `secret_key`)values(? , ? , ?)",
		info.ChannelInfo, info.AppKey, info.SecretKey)
	if err != nil {
		return err
	}
	// save redis
	err = redis.Set(context.TODO(), fmt.Sprintf(redis.SECRET_KEY, info.AppKey), info.SecretKey, 0)
	return err
}

func getSecretInfo(appKey string) (*SecretInfo, error) {
	var res string
	err := redis.Get(context.TODO(), fmt.Sprintf(redis.SECRET_KEY, appKey), &res)

	sInfo := &SecretInfo{}
	if err != nil {
		return nil, err
	} else {
		if res != "" {
			sInfo.SecretKey = res
			sInfo.AppKey = appKey
			return sInfo, nil
		}
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(sqldb.EXPTIME))
	err = sqldb.GetDB(sqldb.MYSQL).Query(ctx, "select * from tomm.channel_infos where app_key = ?", appKey, sInfo)
	cancel()
	if err != nil {
		return nil, err
	}
	return sInfo, nil
}

func createToken(appKey string) string {
	token, _ := getUUID()
	return token
}

func getUUID() (string, error) {

	uuid, err := utils.GetUUID()
	if err != nil {
		return "", nil
	}

	uuidB, _ := uuid.MarshalText()
	res := md5.Sum(uuidB)
	return hex.EncodeToString(res[:]), nil
}
