package redis

import (
	"context"
	"testing"
	"time"
	"tomm/log"
)

func TestExpTime(t *testing.T) {
	for i := 0; i < 10; i++ {
		num := getRandomTime()
		log.Debug("getRandomTime is %d", num)
		time.Sleep(time.Duration(num + 1))
	}
}

func TestRedis(t *testing.T) {
	newRedisCli(nil)
	err := Set(context.Background(), "test", "1111", 0)

	if err != nil {
		log.Error("Redis Set Fail err is %s", err.Error())
		return
	}

}

func TestRedisCmd(t *testing.T) {
	newRedisCli(nil)
	//
	//err := HSet(context.TODO(), "appkey", "test", 11111)
	//
	//if err != nil {
	//	log.Error("Redis Set Fail", zap.String("error", err.Error()))
	//}
	//
	var res string
	err := HGet(context.TODO(), "appkey", "test", &res)

	if err != nil {
		log.Error("Redis Get Fail Err is %s", err.Error())
	}

	//log.Info("Redis Get res is %s", res)
	//var str string
	//err := Get(context.TODO(), "appkey", &str)
	//
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	//err := Exist(context.TODO(), "asda")
	//fmt.Println(err)

	//affRow, err := Del(context.TODO(), "test")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//fmt.Println(affRow)

}
