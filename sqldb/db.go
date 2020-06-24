package sqldb

import (
	"context"
	"sync"
)

type SqlDB interface {
	//getConnStr() string
	Exec(ctx context.Context, sql string, args ...interface{}) error
	Query(ctx context.Context, sql string, res interface{}, args ...interface{}) error
}

type baseConf struct {
}

type DriverType string

var (
	MYSQL DriverType

	mysqlDriver *mysqlDB
	mysqlOnce   *sync.Once
)

func init() {
	mysqlOnce = &sync.Once{}
}

func GetDB(driverType DriverType) SqlDB {
	switch driverType {
	default:
		mysqlOnce.Do(func() {
			mysqlDriver = newMysqlDriver()
		})
		return mysqlDriver
	}
}
