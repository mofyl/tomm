package dao

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"tomm/api/service"
	"tomm/redis"
	"tomm/sqldb"
)

func GetCodeInfoByUserID(mmUserID string) (*service.CodeInfo, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), sqldb.EXPTIME*time.Second)
	codeInfo := &service.CodeInfo{}
	err := sqldb.GetDB(sqldb.MYSQL).QueryOne(ctx, codeInfo, "select * from tomm.code_infos where mm_user_id=?", mmUserID)
	cancel()

	return codeInfo, err
}

func GetCodeInfo(args service.CodeInfo) (*service.CodeInfo, error) {

	ctx, cancel := context.WithTimeout(context.TODO(), sqldb.EXPTIME*time.Second)
	codeInfo := &service.CodeInfo{}
	sqlTotal := strings.Builder{}
	sqlTotal.WriteString("select * from tomm.code_infos where ")
	sql := buildSql(args)
	sqlTotal.WriteString(sql)
	err := sqldb.GetDB(sqldb.MYSQL).QueryOne(ctx, codeInfo, sqlTotal.String())

	cancel()
	return codeInfo, err
}

func CodeExistDB(args service.CodeInfo) (bool, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), sqldb.EXPTIME*time.Second)
	sqlTotal := strings.Builder{}
	sqlTotal.WriteString("select count(*) from tomm.code_infos where ")
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

func buildSql(args service.CodeInfo) string {
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

func SaveCodeInfo(codeInfo *service.CodeInfo, code string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), sqldb.EXPTIME*time.Second)

	if codeInfo.CreateTime == 0 {
		codeInfo.CreateTime = time.Now().Unix()
	}

	_, err := sqldb.GetDB(sqldb.MYSQL).Exec(ctx, "insert into `tomm`.`code_infos`(`app_key`,`create_time`,`mm_user_id`) values(?,?,?)",
		codeInfo.AppKey, codeInfo.CreateTime, codeInfo.MmUserId)
	cancel()
	if err != nil {
		return err
	}

	// 将Code保存到redis
	err = redis.Set(context.TODO(), fmt.Sprintf(redis.CODE_KEY, codeInfo.AppKey, code), codeInfo.MmUserId, redis.CODE_EXP)

	return err

}

func CheckCode(appKey string, code string) (bool, string, error) {
	//key := fmt.Sprintf(redis.CODE_KEY , appKey)
	//redis.Get()

	key := fmt.Sprintf(redis.CODE_KEY, appKey, code)
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
