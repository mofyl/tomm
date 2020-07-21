package dao

import (
	"context"
	"errors"
	"fmt"
	"time"
	"tomm/api/model"
	"tomm/sqldb"
)

const (
	VCODE_EXP = 60
)

func SaveAdminLogin(loginInfo model.AdminInfos) error {

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(sqldb.EXPTIME))
	loginInfo.Created = time.Now().Unix()

	_, err := sqldb.GetDB(sqldb.MYSQL).Exec(ctx, fmt.Sprintf("insert into %s(`login_name`,`login_pwd`,`name`,`number`,`created`) values(?,?,?,?,?) ", ADMIN_INFOS),
		loginInfo.LoginName, loginInfo.LoginPwd, loginInfo.Name, loginInfo.Number, loginInfo.Created)
	cancel()

	if err != nil {
		return err
	}

	// 将名字存在redis中
	SetName(fmt.Sprintf(ADMIN_LOGIN_NAME, loginInfo.LoginName))
	return nil
}

func CheckAdminLoginName(loginName string) (bool, error) {
	return ExistAdminLoginName(loginName)
}

// 判断该用户名是否存在  若存在返回true  不存在返回false
func ExistAdminLoginName(loginName string) (bool, error) {
	return GetName(fmt.Sprintf(ADMIN_LOGIN_NAME, loginName))
}

func GetAdminInfoByLoginName(loginName string) (model.AdminInfos, error) {

	info := model.AdminInfos{}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(sqldb.EXPTIME))

	err := sqldb.GetDB(sqldb.MYSQL).QueryOne(ctx, &info, fmt.Sprintf("select * from %s where login_name=? ", ADMIN_INFOS), loginName)

	cancel()

	return info, err
}

func UpdatePwdByLoginName(loginName string, pwd string) error {

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(sqldb.EXPTIME))

	aff, err := sqldb.GetDB(sqldb.MYSQL).Exec(ctx, fmt.Sprintf("update %s set pwd=? where login_name=? ", ADMIN_INFOS), pwd, loginName)

	cancel()

	if err != nil {
		return err
	}

	rows, err := aff.RowsAffected()

	if err != nil {
		return err
	}

	if rows != 1 {
		return errors.New("Update Fail")
	}

	return nil
}

func DeleteAdminByID(adminID int64) error {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(sqldb.EXPTIME))

	aff, err := sqldb.GetDB(sqldb.MYSQL).Exec(ctx, fmt.Sprintf("delete from %s where id=?", ADMIN_INFOS), adminID)

	cancel()
	if err != nil {
		return err
	}

	rows, err := aff.RowsAffected()

	if err != nil {
		return err
	}

	if rows != 1 {
		return errors.New("Delete Fail")
	}

	return nil

}
