package oauth

import (
	"context"
	"errors"
	"fmt"
	"tomm/redis"
)

const (
	TOKEN_EXP_TIME = 10 * 60 // 10min
)

func GetOAuthInfo(appKey string) (*SecretInfo, error) {
	return getSecretInfo(appKey)
}

func GetToken(appKey string) (string, int64, error) {
	// 查看该Appkey 是否已经存在Token
	var token string
	var err error
	key := fmt.Sprintf(redis.TOKEN_KEY, appKey)
	exist := redis.Exist(context.TODO(), key)

	if exist {
		LeaseRenewKey(appKey, TOKEN_EXP_TIME)
		return token, TOKEN_EXP_TIME, nil
	}
	if token == "" {
		// 表示不存在
		token, err = getUUID()
		if err != nil {
			return "", 0, err
		}
	}
	// 保存到redis中
	err = redis.Set(context.TODO(), fmt.Sprintf(redis.TOKEN_KEY, appKey), token, TOKEN_EXP_TIME)

	if err != nil {
		return "", 0, err
	}

	return token, TOKEN_EXP_TIME, nil
}

func LeaseRenewKey(key string, expTime int64) error {
	// 查看是否存在该key, 不存在直接返回错误
	exist := redis.Exist(context.TODO(), key)

	if !exist {
		return errors.New("Cur Key Not Exist Please Reauthorize")
	}

	// 重新设置为 ex
	res := redis.Expire(context.TODO(), key, expTime)

	if !res {
		return errors.New("Can not Lease Cur Key")
	}
	return nil
}

func CreateOAuthInfo(channelInfo string) (*SecretInfo, error) {
	appKey, err := getUUID()
	if err != nil {
		return nil, err
	}
	secretKey, err := getUUID()

	if err != nil {
		return nil, err
	}

	info := &SecretInfo{
		AppKey:      appKey,
		SecretKey:   secretKey,
		ChannelInfo: channelInfo,
	}

	err = SaveSecretInfo(info)

	if err != nil {
		return nil, err
	}

	return info, nil
}
