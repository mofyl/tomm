package dao

import (
	"context"
	"tomm/redis"
)

func SetName(key string, val interface{}) error {
	_, err := redis.HSet(context.TODO(), NAME_KEY, key, val)
	return err
}

// false 表示不存在 true表示存在
func GetName(key string) (interface{}, error) {

	var val interface{}
	err := redis.HGet(context.TODO(), NAME_KEY, key, &val)

	if err != nil {

		if err != redis.NOT_VALUE {
			return nil, err
		} else {
			return nil, nil
		}

	}

	return val, nil
}

func DelName(key string) (int64, error) {
	return redis.HDel(context.TODO(), NAME_KEY, key)
}
