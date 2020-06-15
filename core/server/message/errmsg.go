package message


type ErrMsg struct {
	ErrCode int `json:"err_code"`
	ErrMsg string `json:"err_msg"`
}

type Msg struct {
	ErrMsg `json:"err_msg"`
	MsgData interface{} `json:"msg_data,omitempty"`
}
