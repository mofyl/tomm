package sqldb

import (
	"context"
	"fmt"
	"testing"
)

func TestMysql(t *testing.T) {
	db := GetDB(MYSQL)
	err := db.Exec(context.TODO(), "select * from channel_infos")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(11111)

}
