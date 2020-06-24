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

type GetTokenRes struct {
	Token   string `json:"token"`
	ExpTime int64  `json:"exp_time"`
}

type BaseMsg struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []byte `json:"data"`
}
