package dao

import (
	"context"
	"errors"
	"fmt"
	"time"
	"tomm/api/model"
	"tomm/log"
	"tomm/redis"
	"tomm/sqldb"
	"tomm/utils"
)

func SavePlatformInfo(info *model.PlatformInfo) error {
	// save DB
	res, err := sqldb.GetDB(sqldb.MYSQL).Exec(context.TODO(),
		fmt.Sprintf("insert into %s(`memo`,`app_key`,`secret_key`,`index_url`,`channel_name`,`sign_url`,`create_time`,`deleted`)values(?,?,?,?,?,?,?,?)", PLATFORM_INFOS),
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
	_, err = redis.HSet(context.TODO(), PLATFORM_HSET, fmt.Sprintf(PLATFORM_INFO_KEY, info.AppKey), infoB)

	// save name
	SetName(fmt.Sprintf(PLATFORM_NAME, info.ChannelName))

	return err
}

func UpdatePlatformInfo(info *model.PlatformInfo) error {

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(sqldb.EXPTIME))
	defer cancel()
	res, err := sqldb.GetDB(sqldb.MYSQL).Exec(ctx, fmt.Sprintf("update %s set channel_info='?' where app_key='?' and deleted=1", PLATFORM_INFOS), info.Memo, info.AppKey)

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

func GetPlatformInfo(appKey string) (*model.PlatformInfo, error) {
	//var res string
	res := &model.PlatformInfo{}
	resB1 := make([]byte, 0)
	key := fmt.Sprintf(PLATFORM_INFO_KEY, appKey)
	err := redis.HGet(context.TODO(), PLATFORM_HSET, key, &resB1)
	//
	if err != nil {
		return nil, err
	}
	//
	err = res.Unmarshal(resB1)

	//sInfo := &api.PlatformInfo{}
	if err != nil {
		return nil, err
	} else if res.AppKey != "" {
		//sInfo.SecretKey = res
		//sInfo.AppKey = appKey
		return res, nil
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(sqldb.EXPTIME))
	err = sqldb.GetDB(sqldb.MYSQL).QueryOne(ctx, res, fmt.Sprintf("select * from %s where app_key='?' and deleted=1", PLATFORM_INFOS), appKey)
	cancel()
	if err != nil {
		return nil, err
	}
	if res.Id == 0 {
		return nil, errors.New("PlatForm Not illegal")
	}

	resB, _ := res.Marshal()
	// 回写到redis中
	_, err = redis.HSet(context.TODO(), PLATFORM_HSET, key, resB)
	if err != nil {
		log.Error("redis Set Fail err is %s", err.Error())
	}
	return res, nil

}

func GetAllPlatform() ([]*model.PlatformInfo, error) {

	//
	//key := fmt.Sprintf(PLATFORM_INFO_KEY, appKey)
	//err := redis.HGet(context.TODO(), PLATFORM_HSET, key, &resB1)

	res, err := redis.HValues(context.TODO(), PLATFORM_HSET)

	if err != nil {
		return nil, err
	}

	infos := make([]*model.PlatformInfo, 0, len(res))

	if len(res) <= 0 {
		// read DB
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(sqldb.EXPTIME))
		err := sqldb.GetDB(sqldb.MYSQL).QueryAll(ctx, &infos, fmt.Sprintf("select * from %s where deleted=1", PLATFORM_INFOS))
		cancel()
		if err != nil {
			return nil, err
		}
		if len(infos) <= 0 {
			return infos, err
		}
		// 写入redis
		resInterface := make([]interface{}, 0, 2*len(infos))

		for i := 0; i < len(infos); i++ {
			b, err := infos[i].Marshal()
			if err != nil {
				continue
			}
			resInterface = append(resInterface, fmt.Sprintf(PLATFORM_INFO_KEY, infos[i].AppKey), b)
		}

		_, err = redis.HSets(context.TODO(), PLATFORM_HSET, resInterface...)

		return infos, err

	} else {
		for i := 0; i < len(res); i++ {
			info := &model.PlatformInfo{}

			err := info.Unmarshal(utils.StrToByte(res[i]))
			if err == nil {
				infos = append(infos, info)
			}
		}
	}

	return infos, nil
}

// true 表示可用 false 表示不可用
func CheckPlatformName(platformName string) bool {

	exist, err := GetName(fmt.Sprintf(PLATFORM_NAME, platformName))

	if err != nil {
		return false
	}
	return !exist
}
