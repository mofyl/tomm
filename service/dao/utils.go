package dao

import (
	"context"
	"tomm/redis"
)

func SetName(key string) error {
	_, err := redis.HSet(context.TODO(), NAME_KEY, key, NAME_VALUE)
	return err
}

// false 表示不存在 true表示存在
func GetName(key string) (bool, error) {

	var val int32
	err := redis.HGet(context.TODO(), NAME_KEY, key, &val)
	if err != nil {
		return false, err
	}

	if val != NAME_VALUE {
		return false, nil
	}

	return true, nil

}
