package jwt

import (
	"go.uber.org/zap"
	"testing"
	"tomm/log"
)

func TestEncode(t *testing.T) {

	j := NewJwt()

	tokenStr, err := j.Encode(TokenInfo{
		ChannelName: "aaaaa",
		ExternInfo:  "",
	})
	if err != nil {
		log.Info("encode", zap.String("Encode Error ", err.Error()))
		return
	}

	log.Info("encode", zap.String("Token ", tokenStr))
	err = j.decode(tokenStr)

	if err != nil {
		log.Info("decode", zap.String("Decode Err", err.Error()))
		return
	}
}
