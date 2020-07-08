package service

import (
	"errors"
	"io/ioutil"
	"tomm/api/service"
	"tomm/ecode"
	"tomm/utils"
)

const (
	GET_LOGINED_USER_INFO = utils.MM_SERVER_URL + "/getUserInfo"
)

type MMRsp struct {
	ErrCode int    `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
	Data    []byte `json:"data"`
}

func GetBaseUserInfo(userID string) (*service.UserBaseInfo, ecode.ErrMsgs) {
	arg := make(map[string]string, 1)
	arg["key"] = utils.MM_PRIVATE_KEY
	arg["user_id"] = userID
	rsp, err := getMMRsp(GET_LOGINED_USER_INFO, arg)
	if err != nil {
		return nil, ecode.MMFail
	}

	if rsp.ErrCode != 1 {
		return nil, ecode.NewErrWithMsg(rsp.ErrMsg, ecode.FromInt(int64(rsp.ErrCode)))
	}

	info := &service.UserBaseInfo{}

	err = utils.Json.Unmarshal(rsp.Data, info)

	return info, ecode.NewErr(err)
}

func getMMRsp(url string, args map[string]string) (*MMRsp, error) {
	rsp, err := HttpGet(url, args)

	if err != nil {
		return nil, err
	}

	if rsp.StatusCode != 200 {
		return nil, errors.New("MM Secret Key Wrong")
	}

	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	mmRsp := &MMRsp{}
	err = utils.Json.Unmarshal(body, mmRsp)

	return mmRsp, err
}
