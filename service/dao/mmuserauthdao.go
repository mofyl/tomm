package dao

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"tomm/api/model"
	"tomm/redis"
	"tomm/sqldb"
)

func GetUserAuthInfoByUserID(mmUserID string) (*model.MMUserAuthInfo, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), sqldb.EXPTIME*time.Second)
	codeInfo := &model.MMUserAuthInfo{}
	err := sqldb.GetDB(sqldb.MYSQL).QueryOne(ctx, codeInfo, fmt.Sprintf("select * from %s where mm_user_id='%s' ", MM_USER_AUTHORIZE_INFO, mmUserID))
	cancel()

	return codeInfo, err
}

func GetUserAuthInfo(args model.MMUserAuthInfo) (*model.MMUserAuthInfo, error) {

	ctx, cancel := context.WithTimeout(context.TODO(), sqldb.EXPTIME*time.Second)
	codeInfo := &model.MMUserAuthInfo{}
	sqlTotal := strings.Builder{}
	sqlTotal.WriteString(fmt.Sprintf("select * from %s where ", MM_USER_AUTHORIZE_INFO))
	sql := buildSql(args)
	sqlTotal.WriteString(sql)
	err := sqldb.GetDB(sqldb.MYSQL).QueryOne(ctx, codeInfo, sqlTotal.String())

	cancel()
	return codeInfo, err
}

func MMUserAuthExistDB(args model.MMUserAuthInfo) (bool, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), sqldb.EXPTIME*time.Second)
	sqlTotal := strings.Builder{}
	sqlTotal.WriteString(fmt.Sprintf("select count(*) from %s where ", MM_USER_AUTHORIZE_INFO))
	sql := buildSql(args)
	sqlTotal.WriteString(sql)
	var res int64
	err := sqldb.GetDB(sqldb.MYSQL).Count(ctx, &res, sqlTotal.String())
	cancel()
	if err != nil {
		return false, err
	}

	return res == 1, nil
}

func buildSql(args model.MMUserAuthInfo) string {
	sql := strings.Builder{}

	if args.Id != 0 {
		sql.WriteString(fmt.Sprintf("id=%d", args.Id))
	}

	if args.MmUserId != "" {
		if sql.Len() != 0 {
			sql.WriteString("and ")
		}
		sql.WriteString(fmt.Sprintf("mm_user_id='%s'", args.MmUserId))
	}

	if args.AppKey != "" {
		if sql.Len() != 0 {
			sql.WriteString("and ")
		}
		sql.WriteString(fmt.Sprintf("app_key='%s'", args.AppKey))
	}
	return sql.String()
}

func SaveMMUserAuthInfo(codeInfo *model.MMUserAuthInfo) error {
	ctx, cancel := context.WithTimeout(context.TODO(), sqldb.EXPTIME*time.Second)

	if codeInfo.CreateTime == 0 {
		codeInfo.CreateTime = time.Now().Unix()
	}

	_, err := sqldb.GetDB(sqldb.MYSQL).Exec(ctx, fmt.Sprintf("insert into %s(`app_key`,`create_time`,`mm_user_id`) values(?,?,?)", MM_USER_AUTHORIZE_INFO),
		codeInfo.AppKey, codeInfo.CreateTime, codeInfo.MmUserId)
	cancel()
	if err != nil {
		return err
	}

	return err

}

func CheckCode(appKey string, code string) (bool, string, error) {
	//key := fmt.Sprintf(redis.CODE_KEY , appKey)
	//redis.Get()

	key := fmt.Sprintf(CODE_KEY, appKey, code)
	var userID string
	err := redis.Get(context.TODO(), key, &userID)

	if err != nil {
		return false, userID, errors.New("Code illegal")
	}

	if userID == "" {
		return false, userID, errors.New("Code illegal")
	}

	aff, err := redis.Del(context.TODO(), key)

	if err != nil {
		return false, userID, errors.New("Code illegal")
	}

	if aff <= 0 {
		return false, userID, errors.New("Code illegal")
	}

	return true, userID, nil
}
