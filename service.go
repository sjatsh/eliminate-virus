package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var sp = "xuehua5201314"

func ModifyUser(openID, cancel string, p, level int64) (*RespData, error) {

	getResultMap := new(RespData)
	getReqMap := make(map[string]interface{})
	getReqMap["plat"] = "wx"
	getReqMap["time"] = time.Now().UnixNano() / 1e6
	getReqMap["openid"] = openID
	getReqMap["wx_appid"] = appId
	getReqMap["wx_secret"] = secret
	getReqMap["sign"] = SignMap(getReqMap)
	delete(getReqMap, "wx_appid")
	delete(getReqMap, "wx_secret")

	getReqData, _ := json.Marshal(getReqMap)

	if err := PostWxGame("/api/archive/get", getReqData, getResultMap); err != nil {
		return nil, err
	}
	if getResultMap.Code != 0 {
		return nil, errors.New("获取用户信息失败")
	}

	recordStr := getResultMap.Data["record"]
	recordMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(fmt.Sprintf("%v", recordStr)), &recordMap); err != nil {
		return nil, err
	}

	resultData, _ := json.Marshal(getResultMap.Data)
	fmt.Printf("当前数据: %s\n", string(resultData))

	// 修改用户数据
	ChangeResult(recordMap, p, level, cancel)

	recordMap["sign"] = SignDataMap(recordMap)
	recordJsonStr, _ := json.Marshal(recordMap)

	reqMap := make(map[string]interface{})
	reqMap["plat"] = "wx"
	reqMap["record"] = string(recordJsonStr)
	reqMap["time"] = time.Now().UnixNano() / 1e6
	reqMap["openid"] = openID
	reqMap["wx_appid"] = appId
	reqMap["wx_secret"] = secret
	reqMap["sign"] = SignMap(reqMap)

	delete(reqMap, "wx_appid")
	delete(reqMap, "wx_secret")
	reqData, _ := json.Marshal(reqMap)

	fmt.Printf("修改后: %s\n", string(reqData))

	uploadResult := new(RespData)
	if err := PostWxGame("/api/archive/upload", reqData, uploadResult); err != nil {
		return nil, err
	}
	return uploadResult, nil
}
