package oauth

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
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
	res, err := sqldb.GetDB(sqldb.MYSQL).Exec(context.TODO(), "insert into channel_infos(`channel_info` , `app_key` , `secret_key`)values(? , ? , ?)",
		info.ChannelInfo, info.AppKey, info.SecretKey)
	if err != nil {
		return err
	}
	if affectNum, err := res.RowsAffected(); err == nil {
		if affectNum <= 0 {
			return errors.New("Insert Fail")
		}
	}

	// save redis
	err = redis.Set(context.TODO(), fmt.Sprintf(redis.SECRET_KEY, info.AppKey), info.SecretKey, 0)
	return err
}

func UpdateChannelInfo(info *SecretInfo) error {

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(sqldb.EXPTIME))
	defer cancel()
	res, err := sqldb.GetDB(sqldb.MYSQL).Exec(ctx, "update tomm.channel_infos set channel_info=? where app_key=?", info.ChannelInfo, info.AppKey)

	if err != nil {
		return err
	}

	if affectNum, err := res.RowsAffected(); err == nil {
		if affectNum <= 0 {
			return errors.New("Update Fail")
		}
	}
	return nil

}

func getSecretInfo(appKey string) (*SecretInfo, error) {
	var res string
	key := fmt.Sprintf(redis.SECRET_KEY, appKey)
	err := redis.Get(context.TODO(), key, &res)

	sInfo := &SecretInfo{}
	if err != nil {
		return nil, err
	} else if res != "" {
		sInfo.SecretKey = res
		sInfo.AppKey = appKey
		return sInfo, nil
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(sqldb.EXPTIME))
	err = sqldb.GetDB(sqldb.MYSQL).Query(ctx, "select * from tomm.channel_infos where app_key = ?", appKey, sInfo)
	cancel()
	if err != nil {
		return nil, err
	}
	// 回写到redis中
	redis.Set(context.TODO(), key, res, 0)
	return sInfo, nil
}

func createToken(appKey string) string {
	token, _ := utils.StrUUID()
	return token
}

//
//func getUUID() (string, error) {
//
//	uuid, err := utils.GetUUID()
//	if err != nil {
//		return "", nil
//	}
//
//	uuidB, _ := uuid.MarshalText()
//	res := md5.Sum(uuidB)
//	return hex.EncodeToString(res[:]), nil
//}
