package dao

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"tomm/api/model"
	"tomm/sqldb"
	"tomm/utils"
)

func GetAllPlatformRole() ([]model.PlatformRoleMidInfo, error) {

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(sqldb.EXPTIME))
	roleInfos := make([]model.PlatformRoleMidInfo, 0)
	err := sqldb.GetDB(sqldb.MYSQL).QueryAll(ctx, &roleInfos, "select a.role_sign,a.create_time as roles_create , a.role_name , b.id as platform_id , b.channel_name ,b.app_key  from `platform`.`platform_infos` as b where b.deleted = 1 RIGHT join `platform`.`platform_roles` as a on a.platform_app_key = b.app_key ")

	cancel()

	if err != nil {
		return nil, err
	}

	return roleInfos, nil

}

func UpdatePlatformRoleName(roleSign, roleName string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(sqldb.EXPTIME))

	aff, err := sqldb.GetDB(sqldb.MYSQL).Exec(ctx, fmt.Sprintf("update %s set role_name=? where role_sign=?", PLATFORM_ROLE), roleName, roleSign)
	cancel()

	if err != nil {
		return 0, err
	}

	return aff.RowsAffected()
}

func DeletePlatformRoleByRoleSign(roleSign string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(sqldb.EXPTIME))
	_, err := sqldb.GetDB(sqldb.MYSQL).Exec(ctx, fmt.Sprintf("delete from %s where role_sign=?", PLATFORM_ROLE), roleSign)
	cancel()

	return err

}

func GetPlatformRoleAppKeyByRoleSign(roleSign string) ([]model.PlatformRole, error) {

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(sqldb.EXPTIME))
	roleInfos := make([]model.PlatformRole, 0)

	err := sqldb.GetDB(sqldb.MYSQL).QueryAll(ctx, &roleInfos, fmt.Sprintf("select role_name , platform_app_key from %s where role_sign=?", PLATFORM_ROLE), roleSign)

	cancel()

	if err != nil {
		return nil, err
	}

	return roleInfos, nil

}

func GetPlatformRoleAppKeyByRoleSigns(roleSign string) ([]model.PlatformRole, error) {

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(sqldb.EXPTIME))
	roleInfos := make([]model.PlatformRole, 0)

	err := sqldb.GetDB(sqldb.MYSQL).QueryAll(ctx, &roleInfos, fmt.Sprintf("select role_name , platform_app_key from %s where role_sign in(?)", PLATFORM_ROLE), roleSign)

	cancel()

	if err != nil {
		return nil, err
	}

	return roleInfos, nil

}

func SavePlatformRole(roleName string, platformAppkeys []string) error {

	nowTime := time.Now().Unix()
	roleID, _ := utils.StrUUID()

	infos := make([]model.PlatformRole, 0, len(platformAppkeys))

	for _, v := range platformAppkeys {

		str := utils.RemoveSpace(v)
		if str != "" {
			info := model.PlatformRole{
				RoleName:       roleName,
				PlatformAppKey: str,
				CreateTime:     nowTime,
				RoleSign:       roleID,
			}

			infos = append(infos, info)
		}

	}

	return savePlatformRole(infos)
}

func savePlatformRole(roleInfo []model.PlatformRole) error {

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("insert into %s(`role_name` , `role_sign` , `platform_app_key` , `create_time`)", PLATFORM_ROLE))

	for i := 0; i < len(roleInfo); i++ {

		builder.WriteString(fmt.Sprintf("values('%s' , '%s' , '%s' , '%d')", roleInfo[i].RoleName, roleInfo[i].RoleSign, roleInfo[i].PlatformAppKey, roleInfo[i].CreateTime))

		if i < len(roleInfo)-1 {
			builder.WriteString(",")
		}
	}

	sqlStr := builder.String()

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(sqldb.EXPTIME))
	_, err := sqldb.GetDB(sqldb.MYSQL).Exec(ctx, sqlStr)
	cancel()

	return err
}

func removePlatformRoleByAppKey(appKey []string) error {

	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("delete from %s where platform_app_key in (", PLATFORM_ROLE))

	for i := 0; i < len(appKey); i++ {

		builder.WriteString(fmt.Sprintf("'%s'", appKey))

		if i < len(appKey)-1 {
			builder.WriteString(",")
		}
	}

	builder.WriteString(")")

	sqlStr := builder.String()
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(sqldb.EXPTIME))
	_, err := sqldb.GetDB(sqldb.MYSQL).Exec(ctx, sqlStr)
	cancel()
	return err
}

func UpdatePlatformRole(roleSign string, roleName string, platformApp map[string]struct{}) error {

	nowTime := time.Now().Unix()

	oldRoleInfo, err := GetPlatformRoleAppKeyByRoleSign(roleSign)

	if err != nil {
		return err
	}

	if len(oldRoleInfo) <= 0 {
		return errors.New("Update Fail : Role Sign illage")
	}

	needUpdateName := false
	oldRoleInfoMap := make(map[string]string, len(oldRoleInfo))

	for i := 0; i < len(oldRoleInfo); i++ {

		appKey := utils.RemoveSpace(oldRoleInfo[i].PlatformAppKey)
		if appKey != "" {
			oldRoleInfoMap[appKey] = oldRoleInfo[i].RoleName
			if oldRoleInfo[i].RoleName != roleName {
				needUpdateName = true
			}
		}
	}

	needAdd := make([]model.PlatformRole, 0, len(oldRoleInfo))
	needDelete := make([]string, 0, len(oldRoleInfo))
	// 现在有 原来没有 就是新加的
	for k, _ := range platformApp {
		name, ok := oldRoleInfoMap[k]
		if !ok {
			info := model.PlatformRole{
				RoleName:       name,
				PlatformAppKey: k,
				CreateTime:     nowTime,
				RoleSign:       roleSign,
			}

			if roleName != "" && roleName != name {
				info.RoleName = roleName
			}

			needAdd = append(needAdd, info)
		}
	}

	// 原来有现在没有 就是需要删除的
	for k, _ := range oldRoleInfoMap {
		_, ok := platformApp[k]

		if !ok {
			needDelete = append(needDelete, k)
		}
	}

	// 删除操作
	if len(needDelete) > 0 {
		err = removePlatformRoleByAppKey(needDelete)
		if err != nil {
			return err
		}
	}

	// 更新之前的名字
	if needUpdateName && len(needDelete) < len(oldRoleInfo) {
		_, err := UpdatePlatformRoleName(roleSign, roleName)

		if err != nil {
			return err
		}

	}

	// 插入操作
	if len(needAdd) > 0 {
		return savePlatformRole(needAdd)
	}
	return nil
}
