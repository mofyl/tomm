package sqldb

import (
	"context"
	"fmt"
	"testing"
	"tomm/api/model"
)

func TestMysql(t *testing.T) {
	//db := GetDB(MYSQL)
	////err := db.Exec(context.TODO(), "select * from platform_infos")
	//
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	//
	//fmt.Println(11111)

	db := GetDB(MYSQL)
	//codeInfo := api.CodeInfo{}
	//err := db.QueryOne(context.TODO(), &codeInfo, "select * from tomm.code_infos where code=?", "qweqwe")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	codeInfos := make([]model.CodeInfo, 0)

	err := db.QueryAll(context.TODO(), &codeInfos, "select * from tomm.code_infos where code=?", "qweqwe")

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range codeInfos {
		fmt.Println(v.CreateTime)
	}

	//res, err := db.Exec(context.TODO(), "insert into tomm.code_infos(`app_key` , `create_time` , `code` , `mm_user_id`) values('1111' , 1594257492, 'qweqwe' , 'asdasd')")
	//
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//affect, err := res.RowsAffected()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(affect)

}
