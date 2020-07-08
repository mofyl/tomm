package dao

import (
	"context"
	"errors"
	"fmt"
	"time"
	"tomm/api/service"
	"tomm/log"
	"tomm/redis"
	"tomm/sqldb"
)

func SavePlatformInfo(info *service.PlatformInfo) error {
	// save DB
	res, err := sqldb.GetDB(sqldb.MYSQL).Exec(context.TODO(),
		"insert into platform_infos(`memo`,`app_key`,`secret_key`,`index_url`,`channel_name`,`sign_url`,`create_time`,`deleted`)values(?,?,?,?,?,?,?,?)",
		info.Memo, info.AppKey, info.SecretKey, info.IndexUrl, info.ChannelName, info.SignUrl, info.CreateTime, info.Deleted)
	if err != nil {
		return err
	}
	if affectNum, err := res.RowsAffected(); err == nil {
		if affectNum <= 0 {
			return errors.New("Insert Fail")
		}
	}
	infoB, _ := info.Marshal()
	// save redis
	err = redis.Set(context.TODO(), fmt.Sprintf(redis.PLATFORM_INFO_KEY, info.AppKey), infoB, -1)
	return err
}

func UpdatePlatformInfo(info *service.PlatformInfo) error {

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(sqldb.EXPTIME))
	defer cancel()
	res, err := sqldb.GetDB(sqldb.MYSQL).Exec(ctx, "update tomm.platform_infos set channel_info=? where app_key=?", info.Memo, info.AppKey)

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

func GetPlatformInfo(appKey string) (*service.PlatformInfo, error) {
	//var res string
	res := &service.PlatformInfo{}
	key := fmt.Sprintf(redis.PLATFORM_INFO_KEY, appKey)
	err := redis.Get(context.TODO(), key, res)

	//sInfo := &service.PlatformInfo{}
	if err != nil {
		return nil, err
	} else if res.AppKey != "" {
		//sInfo.SecretKey = res
		//sInfo.AppKey = appKey
		return res, nil
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(sqldb.EXPTIME))
	err = sqldb.GetDB(sqldb.MYSQL).Query(ctx, res, "select * from tomm.platform_infos where app_key = ?", appKey)
	cancel()
	if err != nil {
		return nil, err
	}
	resB, _ := res.Marshal()
	// 回写到redis中
	redis.Set(context.TODO(), key, resB, -1)
	return res, nil

}

// true 表示可用 false 表示不可用
func CheckPlatformName(platformName string) bool {

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(sqldb.EXPTIME))
	res := &service.PlatformInfo{}
	err := sqldb.GetDB(sqldb.MYSQL).Query(ctx, res, "select channel_name from tomm.platform_infos where channel_name = ?", platformName)
	cancel()
	if err != nil {
		log.Error("Check PlatformName Fail Error is %s", err.Error())
		return false
	}

	if res.ChannelName == "" {
		return true
	}

	return false

}
