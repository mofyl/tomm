package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type redisConf struct {
	NetWork      string
	Addr         string
	UserName     string
	Pwd          string
	DB           int
	IdleTimeout  time.Duration // default 5min
	DialTimeout  time.Duration // default 5s
	WriteTimeout time.Duration // default 3s
	ReadTimeout  time.Duration // default 3s
	//MaxLifeTime  time.Duration // 默认不给这个字段 如果后面需要可以添加
}

var (
	cli *redis.Ring
)

func init() {

}

func defaultConf() *redisConf {
	return &redisConf{
		NetWork:      "",
		Addr:         "",
		UserName:     "",
		Pwd:          "",
		DB:           0,
		IdleTimeout:  0,
		DialTimeout:  0,
		WriteTimeout: 0,
		ReadTimeout:  0,
	}
}

func newRedisCli(conf *redisConf) *redis.Ring {
	if conf == nil {
		conf = defaultConf()
	}
	opt := &redis.RingOptions{
		Username:     conf.UserName,
		Password:     conf.Pwd,
		DB:           conf.DB,
		IdleTimeout:  conf.IdleTimeout,
		DialTimeout:  conf.DialTimeout,
		WriteTimeout: conf.WriteTimeout,
		ReadTimeout:  conf.ReadTimeout,
	}
	cli := redis.NewRing(opt)
	ctx, cancel := context.WithTimeout(context.Background(), conf.IdleTimeout)
	defer cancel()
	if err := cli.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Ping Redis Fail addr is %s", conf.Addr))
	}
	return cli
}

func Get(ctx context.Context, key string, data interface{}) error {
	res := cli.Get(ctx, key)

	if res.Err() != nil {
		return res.Err()
	}

	return res.Scan(data)
}

func Del(ctx context.Context, key string) (int64, error) {
	count, err := cli.Del(ctx, key).Result()

	if err != nil {
		return 0, err
	} else {
		return count, err
	}
}
