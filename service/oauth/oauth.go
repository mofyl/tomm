package oauth

import (
	"context"
	"errors"
	"fmt"
	"tomm/redis"
)

const (
	TOKEN_EXP_TIME = 2 * 3600
)

func GetOAuthInfo(appKey string) (*SecretInfo, error) {
	return getSecretInfo(appKey)
}

func GetToken(appKey string) (string, error) {
	// 查看该Appkey 是否已经存在Token
	var token string
	err := redis.Get(context.TODO(), fmt.Sprintf(redis.SECRET_KEY, appKey), &token)
	if err != nil {
		return "", err
	}

	if token == "" {
		// 表示不存在
		token, err = getUUID()
		if err != nil {
			return "", err
		}
	}

	// 保存到redis中
	err = redis.Set(context.TODO(), fmt.Sprintf(redis.TOKEN_KEY, appKey), token, TOKEN_EXP_TIME)

	if err != nil {
		return "", nil
	}

	return token, nil
}

func LeaseRenewToken(appKey string) error {
	// 查看是否存在该key, 不存在直接返回错误
	key := fmt.Sprintf(redis.TOKEN_KEY, appKey)
	exist := redis.Exist(context.TODO(), key)

	if !exist {
		return errors.New("Cur Key Not Exist Please Reauthorize")
	}

	// 重新设置为 ex
	res := redis.Expire(context.TODO(), key, TOKEN_EXP_TIME)

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
