package service

import (
	"errors"
	"io/ioutil"
	"tomm/api/model"
	"tomm/ecode"
	"tomm/log"
	"tomm/utils"
)

var mmSerUrl string

type MMRsp struct {
	ErrCode int    `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
	Data    []byte `json:"data"`
}

func GetBaseUserInfo(userID string) (*model.UserBaseInfo, ecode.ErrMsgs) {
	arg := make(map[string]string, 1)
	arg["key"] = utils.MM_PRIVATE_KEY
	arg["user_id"] = userID
	rsp, err := getMMRsp(mmSerUrl+"/getUserInfo", arg)
	if err != nil {
		log.Error("GetBaseUserInfo Err is %s", err.Error())
		return nil, ecode.MMFail
	}

	if rsp.ErrCode != 1 {
		return nil, ecode.NewErrWithMsg(rsp.ErrMsg, ecode.FromInt(int64(rsp.ErrCode)))
	}

	info := &model.UserBaseInfo{}

	err = utils.Json.Unmarshal(rsp.Data, info)
	if err != nil {
		return nil, ecode.NewErr(err)
	}
	return info, nil
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
