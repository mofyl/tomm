package dao

import (
	"context"
	"fmt"
	"strings"
	"time"
	"tomm/api/service"
	"tomm/sqldb"
)

func GetCodeInfoByUserID(mmUserID string) (*service.CodeInfo, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), sqldb.EXPTIME*time.Second)
	codeInfo := &service.CodeInfo{}
	err := sqldb.GetDB(sqldb.MYSQL).Query(ctx, codeInfo, "select * from tomm.code_infos where mm_user_id=?", mmUserID)
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
	err := sqldb.GetDB(sqldb.MYSQL).Query(ctx, codeInfo, sqlTotal.String())

	cancel()
	return codeInfo, err
}

func CodeExist(args service.CodeInfo) (bool, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), sqldb.EXPTIME*time.Second)
	sqlTotal := strings.Builder{}
	sqlTotal.WriteString("select count(*) from tomm.code_infos where ")
	sql := buildSql(args)
	sqlTotal.WriteString(sql)
	res, err := sqldb.GetDB(sqldb.MYSQL).Exec(ctx, sqlTotal.String())
	cancel()
	if err != nil {
		return false, err
	}

	affect, err := res.RowsAffected()

	if err != nil {
		return false, err
	}

	if affect >= 0 {
		return false, nil
	}

	return true, nil

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
		sql.WriteString(fmt.Sprintf("mm_user_id=%s", args.MmUserId))
	}

	if args.Code != "" {
		if sql.Len() != 0 {
			sql.WriteString("and ")
		}
		sql.WriteString(fmt.Sprintf("code=%s", args.Code))
	}

	if args.AppKey != "" {
		if sql.Len() != 0 {
			sql.WriteString("and ")
		}
		sql.WriteString(fmt.Sprintf("app_key=%s", args.AppKey))
	}
	return sql.String()
}

func SaveCodeInfo(codeInfo *service.CodeInfo) error {
	ctx, cancel := context.WithTimeout(context.TODO(), sqldb.EXPTIME*time.Second)

	if codeInfo.CreateTime == 0 {
		codeInfo.CreateTime = time.Now().Unix()
	}

	_, err := sqldb.GetDB(sqldb.MYSQL).Exec(ctx, "insert into `tomm`.`code_infos`(`app_key`,`create_time`,`code`,`mm_user_id`) values(?,?,?,?)",
		codeInfo.AppKey, codeInfo.CreateTime, codeInfo.Code, codeInfo.MmUserId)
	cancel()
	if err != nil {
		return err
	}
	return nil
	//
	//if res.RowsAffected() == int64(0) {
	//	return errors.New("")
	//}

}
