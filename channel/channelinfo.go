package channel

import (
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

	for _, v := range channelMap {

		log.Info("Channel Info %v ", v)

	}

}

func init() {

	channelMap = make(map[string]*ChannelInfo)

	err := config.DecodeAll(CHANNEL_FILE_NAME, &channelMap)

	if err != nil {
		panic("Init Channel Info Fail " + err.Error())
		return
	}

	log.Debug("Init ChannelInfo")
}

func GetChannelInfo(channelName string) *ChannelInfo {
	info, ok := channelMap[channelName]

	if !ok {
		return nil
	}

	return info
}
