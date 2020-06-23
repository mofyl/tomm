package channel

import (
	"go.uber.org/zap"
	"tomm/config"
	"tomm/log"
)

type ChannelInfo struct {
	ChannelID   string `yaml:"channelID"`
	ChannelName string `yaml:"channelName"`
	LoginUrl    string `yaml:"loginUrl"`
}

const (
	CHANNEL_FILE_NAME = "channel_info.yaml"
)

var (
	channelMap map[string]*ChannelInfo // channel_name
)

func testChannelInfo() {

	channelMap = make(map[string]*ChannelInfo)

	err := config.DecodeAll(CHANNEL_FILE_NAME, &channelMap)

	if err != nil {
		panic("Init Channel Info Fail " + err.Error())
		return
	}

	for k, v := range channelMap {

		log.Msg(log.DEBUG, "Channel Name is  "+k)
		log.Info("Channel Info", zap.Any("ChannelInfo ", v))

	}

}

func init() {

	channelMap = make(map[string]*ChannelInfo)

	err := config.DecodeAll(CHANNEL_FILE_NAME, &channelMap)

	if err != nil {
		panic("Init Channel Info Fail " + err.Error())
		return
	}

	log.Msg(log.DEBUG, "Init ChannelInfo")
}

func GetChannelInfo(channelName string) *ChannelInfo {
	info, ok := channelMap[channelName]

	if !ok {
		return nil
	}

	return info
}
