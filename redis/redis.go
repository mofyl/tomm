package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"math/rand"
	"time"
	"tomm/config"
	"tomm/log"
)

type RedisConf struct {
	Addr         string `yaml:"addr"`
	UserName     string `yaml:"userName"`
	Pwd          string `yaml:"pwd"`
	DB           int    `yaml:"sqldb"`
	IdleTimeout  int64  `yaml:"idleTimeout"`  // default 5min
	DialTimeout  int64  `yaml:"dialTimeout"`  // default 5s
	WriteTimeout int64  `yaml:"writeTimeout"` // default 3s
	ReadTimeout  int64  `yaml:"readTimeout"`  // default 3s
	MinExpTime   int64  `yaml:"minExpTime"`
	//MaxLifeTime  time.Duration // 默认不给这个字段 如果后面需要可以添加
}

var (
	//cli *redis.Ring
	cli         *redis.Client
	defaultConf *RedisConf
)

func init() {
	defaultConf = &RedisConf{}
	if err := config.Decode(config.CONFIG_FILE_NAME, "redis", defaultConf); err != nil {
		panic("Decode Redis Conf Fail " + err.Error())
	}
}

func newRedisCli(conf *RedisConf) {
	if conf == nil {
		conf = defaultConf
	}
	opt := &redis.Options{
		Username:     conf.UserName,
		Password:     conf.Pwd,
		DB:           conf.DB,
		IdleTimeout:  time.Duration(conf.IdleTimeout) * time.Second,
		DialTimeout:  time.Duration(conf.DialTimeout) * time.Second,
		WriteTimeout: time.Duration(conf.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Second,
	}
	cli = redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.DialTimeout)*time.Second)
	defer cancel()
	if err := cli.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Ping Redis Fail addr is %s ,errmsg is %s\n", conf.Addr, err.Error()))
	}
}

func newRedisRingCli(conf *RedisConf) *redis.Ring {
	if conf == nil {
		conf = defaultConf
	}
	opt := &redis.RingOptions{
		Username:     conf.UserName,
		Password:     conf.Pwd,
		DB:           conf.DB,
		IdleTimeout:  time.Duration(conf.IdleTimeout) * time.Second,
		DialTimeout:  time.Duration(conf.DialTimeout) * time.Second,
		WriteTimeout: time.Duration(conf.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Second,
	}
	cli := redis.NewRing(opt)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.DialTimeout)*time.Second)
	defer cancel()
	if err := cli.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Ping Redis Fail addr is %s ,errmsg is %s\n", conf.Addr, err.Error()))
	}
	return cli
}

// 这里若expTime给 0 表示该key永远不过期
// 这里若给-1 表示设置随机值
func Set(ctx context.Context, key string, data interface{}, expTime int64) error {
	if expTime == -1 {
		// 给随机值
		expTime = defaultConf.MinExpTime + getRandomTime()
		log.Debug("ExpTime is ", zap.Int64("ExpTime ", expTime))
	}
	res := cli.Set(ctx, key, data, time.Duration(expTime)*time.Second)
	if res.Err() != nil {
		return res.Err()
	}

	if res.Val() != "OK" {
		return errors.New("Redis Set Key Fail Key is " + key)
	}

	return nil
}

func Exist(ctx context.Context, key ...string) bool {
	cmd := cli.Exists(ctx, key...)

	if cmd.Err() != nil {
		return false
	}

	if cmd.Val() != 1 {
		return false
	}
	return true
}

func Expire(ctx context.Context, key string, expTime int64) bool {
	cmd := cli.Expire(ctx, key, time.Duration(expTime)*time.Second)

	return cmd.Val()
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

func HSet(ctx context.Context, key string, field string, value interface{}) error {
	cmd := cli.HSet(ctx, key, field, value)
	return cmd.Err()
}

func HGet(ctx context.Context, key string, field string, value interface{}) error {
	cmd := cli.HGet(ctx, key, field)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return cmd.Scan(value)
}

func getRandomTime() int64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(defaultConf.MinExpTime)
}
