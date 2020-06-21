package sqldb

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
	"tomm/config"
)

type mysqlConf struct {
	baseConf
	MaxOpenConn  int `yaml:"maxOpenConn"`
	MaxIdleConn  int `yaml:"maxIdleConn"`
	ConnLiftTime int `yaml:"connLiftTime"`
}

type mysqlDB struct {
	engine *sqlx.DB
	conf   *mysqlConf
}

func newMysqlDriver() *mysqlDB {

	db := mysqlDB{}
	conf := &mysqlConf{}
	err := config.Decode("mysql", conf)
	if err != nil {
		panic("Create Mysql Driver Fail " + err.Error())
	}
	db.conf = conf
	var engine *sqlx.DB
	engine, err = db.Connect()
	if err != nil {
		panic("Mysql Connect Fail " + err.Error())
	}

	engine.SetConnMaxLifetime(time.Duration(db.conf.ConnLiftTime) * time.Second)
	engine.SetMaxIdleConns(db.conf.MaxIdleConn)
	engine.SetMaxOpenConns(db.conf.MaxOpenConn)

	db.engine = engine

	return nil
}

func (m *mysqlDB) Connect() (*sqlx.DB, error) {
	connStr := m.getConnStr()
	db, err := sqlx.Connect("mysql", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (m *mysqlDB) getConnStr() string {
	return fmt.Sprintf("%s:%s@tcp(%s)?charset=utf8", m.conf.UserName, m.conf.UserName, m.conf.Addr)
}

func (m *mysqlDB) Query() {
}
