package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

type RespData struct {
	Data map[string]interface{} `json:"data"`
	Code int                    `json:"code"`
}

func PostWxGame(uri string, req []byte, respModel interface{}) error {
	resp, err := http.Post(gameUrl+uri, "application/json", bytes.NewReader(req))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if respModel != nil {
		if err := json.Unmarshal(body, respModel); err != nil {
			return err
		}
	}
	return nil
}

func SignMap(dataMap map[string]interface{}) string {
	reqKeys := make([]string, 0, len(dataMap))
	for key := range dataMap {
		reqKeys = append(reqKeys, key)
	}
	sort.Strings(reqKeys)

	reqSignEle := make([]string, 0)
	for i, j := 0, len(reqKeys); i < j; i++ {
		reqSignEle = append(reqSignEle, reqKeys[i]+"="+fmt.Sprintf("%v", dataMap[reqKeys[i]]))
	}
	reqBS := md5.Sum(bytes.NewBufferString(strings.Join(reqSignEle, "&")).Bytes())
	return strings.ToLower(hex.EncodeToString(reqBS[:]))
}

func SignDataMap(dataMap map[string]interface{}) string {
	// 取出所有的键，并排序
	keys := make([]string, 0, len(dataMap))
	for key := range dataMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// 拼接参数
	ignS := ""
	for i, j := 0, len(keys); i < j; i++ {
		ignS += fmt.Sprintf("%v", dataMap[keys[i]])
	}
	bs := md5.Sum(bytes.NewBufferString(ignS).Bytes())
	return strings.ToLower(hex.EncodeToString(bs[:]))
}
