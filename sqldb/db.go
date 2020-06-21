package sqldb

import "sync"

type SqlDB interface {
	//getConnStr() string
}

type baseConf struct {
	Addr     string `yaml:"addr"`
	UserName string `yaml:"userName"`
	Pwd      string `yaml:"pwd"`
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
