package oauth

import (
	"context"
	"errors"
	"fmt"
	"tomm/redis"
	"tomm/utils"
)

const (
	RESOURCE_TOKEN_EXP = 10 * 60 // 10min
)

func GetOAuthInfo(appKey string) (*SecretInfo, error) {
	return getSecretInfo(appKey)
}

func GetToken(appKey string) (string, int64, error) {
	// 查看该Appkey 是否已经存在Token
	var token string
	var err error
	key := fmt.Sprintf(redis.RESOURCE_KEY, appKey)
	err = redis.Get(context.TODO(), key, &token)
	//exist := redis.Exist(context.TODO(), key)

	if token != "" && err == nil {
		LeaseRenewKey(key, RESOURCE_TOKEN_EXP)
		return token, RESOURCE_TOKEN_EXP, nil
	}
	if token == "" {
		// 表示不存在
		token, err = utils.StrUUID()
		if err != nil {
			return "", 0, err
		}
	}
	// 保存到redis中
	err = redis.Set(context.TODO(), fmt.Sprintf(redis.RESOURCE_KEY, appKey), token, RESOURCE_TOKEN_EXP)

	if err != nil {
		return "", 0, err
	}

	return token, RESOURCE_TOKEN_EXP, nil
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
	appKey, err := utils.StrUUID()
	if err != nil {
		return nil, err
	}
	secretKey, err := utils.StrUUID()

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
