package oauth

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"tomm/sqldb"
	"tomm/utils"
)

type SecretInfo struct {
	ChannelInfo string
	SecretKey   string
	AppKey      string
}

func getSecretInfo(channelInfo string) (*SecretInfo, error) {
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

	// save DB
	err = sqldb.GetDB(sqldb.MYSQL).Exec(context.TODO(), "insert into channel_infos(`channel_info` , `app_key` , `secret_key`)values(? , ? , ?)",
		info.ChannelInfo, info.AppKey, info.SecretKey)
	// save redis
	// 设置到一个 HSet中

	return info, err
}

func getUUID() (string, error) {

	uuid, err := utils.GetUUID()
	if err != nil {
		return "", nil
	}

	uuidB, _ := uuid.MarshalText()
	res := md5.Sum(uuidB)
	return hex.EncodeToString(res[:]), nil
}
