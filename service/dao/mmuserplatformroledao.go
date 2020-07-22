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

func GetMMUserPlatformRoleSign(userID string) ([]model.MMUserPlatformRoles, error) {

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(sqldb.EXPTIME))
	mmUserRoleInfo := make([]model.MMUserPlatformRoles, 0)
	err := sqldb.GetDB(sqldb.MYSQL).QueryAll(ctx, &mmUserRoleInfo, fmt.Sprintf("select role_sign  from %s where mm_user_id=?", MM_USER_PLATFORM_ROLE), userID)
	cancel()

	if err != nil {
		return nil, err
	}

	return mmUserRoleInfo, nil

}

func addMMUserPlatformRole(info []model.MMUserPlatformRoles) (int64, error) {

	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("insert into %s(`mm_user_id`,`role_sign`) ", MM_USER_PLATFORM_ROLE))

	for i := 0; i < len(info); i++ {

		builder.WriteString(fmt.Sprintf("values('%s' , '%s')", info[i].MmUserId, info[i].RoleSign))

		if i <= len(info)-1 {
			builder.WriteString(", ")
		}
	}

	sqlStr := builder.String()
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(sqldb.EXPTIME))
	row, err := sqldb.GetDB(sqldb.MYSQL).Exec(ctx, sqlStr)
	cancel()

	if err != nil {
		return 0, err
	}

	return row.RowsAffected()
}

func AddMMUserPlatformRole(mmUserID string, roleSign []string) error {

	infos := make([]model.MMUserPlatformRoles, 0, len(roleSign))

	for i := 0; i < len(roleSign); i++ {
		info := model.MMUserPlatformRoles{
			MmUserId: mmUserID,
			RoleSign: roleSign[i],
		}

		infos = append(infos, info)
	}

	aff, err := addMMUserPlatformRole(infos)

	if err != nil {
		return nil
	}

	if aff <= 0 {
		return errors.New("Affect Rows is Zero")
	}

	return nil

}
