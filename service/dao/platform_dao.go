package dao

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"tomm/api/model"
	"tomm/sqldb"
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
	curId, err := res.LastInsertId()
	if err != nil {
		return err
	}
	//infoB, _ := info.Marshal()
	// save redis
	//_, err = redis.HSet(context.TODO(), PLATFORM_HSET, fmt.Sprintf(PLATFORM_INFO_KEY, info.AppKey), infoB)

	// save name
	SetName(fmt.Sprintf(PLATFORM_NAME, info.ChannelName), curId)

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

//
//func DeletePlatformInfo(appKey string) error {
//	// TODO：这里考虑要不要做成事务
//
//	// 除了删除平台信息外 还将权限信息都要删除
//
//}

func GetPlatformInfo(appKey string) (*model.PlatformInfo, error) {
	//var res string
	res := &model.PlatformInfo{}
	//resB1 := make([]byte, 0)
	//key := fmt.Sprintf(PLATFORM_INFO_KEY, appKey)
	//err := redis.HGet(context.TODO(), PLATFORM_HSET, key, &resB1)
	//
	//if err != nil {
	//	return nil, err
	//}
	//
	//err = res.Unmarshal(resB1)

	//sInfo := &api.PlatformInfo{}
	//if err != nil {
	//	return nil, err
	//} else if res.AppKey != "" {
	//	//sInfo.SecretKey = res
	//	//sInfo.AppKey = appKey
	//	return res, nil
	//}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(sqldb.EXPTIME))
	err := sqldb.GetDB(sqldb.MYSQL).QueryOne(ctx, res, fmt.Sprintf("select id,memo,app_key,index_url,channel_name,sign_url,create_time from %s where app_key='?' and deleted=1", PLATFORM_INFOS), appKey)
	cancel()
	if err != nil {
		return nil, err
	}
	if res.Id == 0 {
		return nil, errors.New("PlatForm illegal")
	}

	//resB, _ := res.Marshal()
	// 回写到redis中
	//_, err = redis.HSet(context.TODO(), PLATFORM_HSET, key, resB)
	//if err != nil {
	//	log.Error("redis Set Fail err is %s", err.Error())
	//}
	return res, nil

}

func GetPlatformByAppKeys(appKeys map[string]struct{}) ([]*model.PlatformInfo, error) {

	if appKeys == nil || len(appKeys) <= 0 {
		return nil, nil
	}

	builder := strings.Builder{}
	//appKeyStr.WriteString("(")
	for k, _ := range appKeys {
		builder.WriteString("'")
		builder.WriteString(k)
		builder.WriteString("',")
	}

	//appKeyStr.WriteString(")")

	appKeyStr := builder.String()
	appKeyStr = appKeyStr[:len(appKeyStr)-1]

	infos := make([]*model.PlatformInfo, 0, len(appKeys))

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(sqldb.EXPTIME))
	err := sqldb.GetDB(sqldb.MYSQL).QueryAll(ctx, &infos, fmt.Sprintf("select id,memo,app_key,index_url,channel_name,sign_url,create_time from %s where app_key in (%s) and deleted=1", PLATFORM_INFOS, appKeyStr))
	cancel()
	if err != nil {
		return nil, err
	}
	if len(infos) <= 0 {
		return nil, errors.New("PlatForm Not illegal")
	}

	return infos, nil

}

func GetPlatformCount() (int64, error) {

	var count int64

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(sqldb.EXPTIME))

	err := sqldb.GetDB(sqldb.MYSQL).Count(ctx, &count, fmt.Sprintf("select count(id) from %s where deleted=1", PLATFORM_INFOS))
	cancel()

	if err != nil {
		return 0, err
	}

	return count, nil

}

func DeletePlatformByNames(names []string) error {

	ids := strings.Builder{}

	for _, v := range names {

		val, _ := GetName(v)
		if val != nil {
			if valI, ok := val.(string); ok {
				ids.WriteString(valI)
				ids.WriteString(",")
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(sqldb.EXPTIME))
	if ids.Len() > 0 {
		builder := ids.String()
		builder = builder[:len(builder)-1]
		_, err := sqldb.GetDB(sqldb.MYSQL).Exec(ctx, fmt.Sprintf("update %s set deleted=2  where id in(?) and deleted=1", PLATFORM_INFOS), builder)
		cancel()

		if err != nil {
			return err
		}
	} else if ids.Len() == 0 && len(names) > 0 {
		// 拼接 names
		build := strings.Builder{}
		build.WriteString(fmt.Sprintf("update %s set deleted=2  where channel_name in(", PLATFORM_INFOS))
		for i := 0; i < len(names); i++ {
			build.WriteString(fmt.Sprintf("'%s'", names[i]))

			if i < len(names)-1 {
				build.WriteString(",")
			}
		}

		build.WriteString(")and deleted=1 ")
		sqlBuild := build.String()
		_, err := sqldb.GetDB(sqldb.MYSQL).Exec(ctx, sqlBuild)
		cancel()

		if err != nil {
			return err
		}
	}

	for _, v := range names {
		DelName(fmt.Sprintf(PLATFORM_NAME, v))
	}

	return nil
}

func GetAllPlatform(page, pageSize int32) ([]*model.PlatformInfo, error) {

	//
	//key := fmt.Sprintf(PLATFORM_INFO_KEY, appKey)
	//err := redis.HGet(context.TODO(), PLATFORM_HSET, key, &resB1)

	//res, err := redis.HValues(context.TODO(), PLATFORM_HSET)
	//
	//if err != nil {
	//	return nil, err
	//}

	//infos := make([]*model.PlatformInfo, 0, len(res))

	infos := make([]*model.PlatformInfo, 0, pageSize)
	//if len(res) >{
	// read DB
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(sqldb.EXPTIME))
	err := sqldb.GetDB(sqldb.MYSQL).QueryAll(ctx, &infos, fmt.Sprintf("select id,memo,app_key,index_url,channel_name,sign_url,create_time from %s where deleted=1 limit ?,?", PLATFORM_INFOS), page*pageSize, pageSize)
	cancel()
	if err != nil {
		return nil, err
	}
	//if len(infos) <= 0 {
	//	return infos, err
	//}
	// 写入redis
	//resInterface := make([]interface{}, 0, 2*len(infos))
	//
	//for i := 0; i < len(infos); i++ {
	//	b, err := infos[i].Marshal()
	//	if err != nil {
	//		continue
	//	}
	//	resInterface = append(resInterface, fmt.Sprintf(PLATFORM_INFO_KEY, infos[i].AppKey), b)
	//}
	//
	//_, err = redis.HSets(context.TODO(), PLATFORM_HSET, resInterface...)

	return infos, err

	//} else {
	//	for i := 0; i < len(res); i++ {
	//		info := &model.PlatformInfo{}
	//
	//		err := info.Unmarshal(utils.StrToByte(res[i]))
	//		if err == nil {
	//			infos = append(infos, info)
	//		}
	//	}
	//}

	//return infos, nil
}

// true 表示可用 false 表示不可用
func CheckPlatformName(platformName string) bool {

	res, err := GetName(fmt.Sprintf(PLATFORM_NAME, platformName))

	if err != nil {
		return false
	}
	if res == nil {
		return true
	}

	return false
}
