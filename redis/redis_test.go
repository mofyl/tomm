package redis

import (
	"context"
	"go.uber.org/zap"
	"testing"
	"time"
	"tomm/log"
)

func TestExpTime(t *testing.T) {
	for i := 0; i < 10; i++ {
		num := getRandomTime()
		log.Debug("getRandomTime", zap.Int64("randomTime", num))
		time.Sleep(time.Duration(num + 1))
	}
}

func TestRedis(t *testing.T) {
	newRedisCli(nil)
	err := Set(context.Background(), "test", "1111", 0)

	if err != nil {
		log.Error("Redis Set Fail", zap.String("error", err.Error()))
		return
	}
	log.Msg(log.DEBUG, "Success")

}

func TestRedisCmd(t *testing.T) {
	newRedisCli(nil)
	//
	//err := HSet(context.TODO(), "appkey", "test", 11111)
	//
	//if err != nil {
	//	log.Error("Redis Set Fail", zap.String("error", err.Error()))
	//}

	var res string
	err := HGet(context.TODO(), "appkey", "test", &res)

	if err != nil {
		log.Error("Redis Get Fail", zap.String("error", err.Error()))
	}

	log.Info("Redis Get ", zap.String("res", res))
}
