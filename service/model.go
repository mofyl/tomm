package service

const (
	DATALEN          = 4
	TIMELEN          = 4
	CHANNEL_INFO_LEN = 4
	MAX_DATA         = 512

	MAX_TTL = 10 * 60 // 若数据包: nowTime - sendTime > 10min 则不处理
)

type GetTokenReq struct {
	AppKey string `form:"app_key"`
	Data   string `form:"data"`
}

type ReqDataInfo struct {
	SendTime    int64
	ChannelInfo string
	ExtendInfo  []byte
	// 两个url
	Resource    string // 这里的资源其实就是url
	Method      string
	ResourceUrl string // 表示资源服务器的资源  这里的资源服务器一定要实现 /verify
	BackUrl     string // 第三方回调地址
}

type GetTokenRes struct {
	TokenInfo string `json:"token_info,omitempty"`
}

type TokenInfo struct {
	Token      string `json:"token,omitempty"`
	ExpTime    int64  `json:"exp_time,omitempty"`
	ExtendInfo []byte `json:"extend_info,omitempty"`
}

type VerifyTokenReq struct {
	AppKey string `form:"app_key"`
	Token  string `form:"token"`
}

type VerifyTokenRes struct {
	ExpTime int64 `json:"exp_time,omitempty"`
}
