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
	DataLen     int32
	SendTime    int64
	ChannelInfo string
	ExtendInfo  []byte
}

type GetTokenRes struct {
	Token      string `json:"token"`
	ExpTime    int64  `json:"exp_time"`
	ExtendInfo []byte `json:"extern_info"`
	ErrCode    int64  `json:"err_code"`
	ErrMsg     string `json:"err_msg"`
}

type VerifyTokenReq struct {
	AppKey string `form:"app_key"`
	Token  string `form:"token"`
}

type VerifyTokenRes struct {
	ExpTime int64  `json:"exp_time"`
	ErrCode int64  `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}
