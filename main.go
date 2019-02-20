package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

var secret = "8fbd540d0b23197df1d5095f0d6ee46d"
var appId = "wxa2c324b63b2a9e5e"
var gameUrl = "https://wxwyjh.chiji-h5.com"
var m int64 = 1000000

type RespData struct {
	Data map[string]interface{} `json:"data"`
	Code int                    `json:"code"`
}

func main() {

	p := flag.Int64("p", 0, "选择套餐")
	level := flag.Int("level", 0, "关卡等级")
	openID := flag.String("id", "", "open_id")
	flag.Parse()

	if "" == *openID {
		log.Fatal("请填写用户id")
	}
	if *p == 8 && *level <= 0 {
		log.Fatal("请填写关卡等级")
	}
	getResultMap := new(RespData)

	getReqMap := make(map[string]interface{})
	getReqMap["plat"] = "wx"
	getReqMap["time"] = time.Now().UnixNano() / 1e6
	getReqMap["openid"] = *openID
	getReqMap["wx_appid"] = appId
	getReqMap["wx_secret"] = secret
	getReqMap["sign"] = SignMap(getReqMap)
	delete(getReqMap, "wx_appid")
	delete(getReqMap, "wx_secret")

	getReqData, _ := json.Marshal(getReqMap)

	if err := PostWxGame("/api/archive/get", getReqData, getResultMap); err != nil {
		log.Fatal(err)
	}
	if getResultMap.Code != 0 {
		log.Fatal("获取用户信息失败")
	}

	recordStr := getResultMap.Data["record"]
	recordMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(fmt.Sprintf("%v", recordStr)), &recordMap); err != nil {
		log.Fatal(err)
	}

	if *p <= 0 {
		log.Fatal("请选择对应套餐")
	}

	switch *p {
	//套餐1: 主武器满级+金币收益满级18元
	case 1:
		recordMap["lDamage"] = 999
		recordMap["lCount"] = 356
		recordMap["lJiaZhi"] = 999
		recordMap["lRiChang"] = 999
	case 2:
		// 套餐2: 七个副武器满级打包30元
		recordMap["levelFuCount"] = "[32,32,32,32,32,32,32,1,1,1]"
		recordMap["levelFuDamage"] = "[999,999,999,999,999,999,999,1,1,1]"
	case 3:
		// 套餐3: 无限体力18元
		recordMap["tiLi"] = 999 * m
	case 4:
		// 套餐4: 无线钻石18元（需要自己一个个兑换金币）
		recordMap["zuanShi"] = 999999999 * m
	case 5:
		// 套餐5: 无限金币24元（武器全部升级满级用不完）
		recordMap["money"] = 999999999 * m
	case 6:
		// 套餐5: 无限金币、钻石、体力（可以自己升级体验游戏乐趣）
		recordMap["money"] = 999999999 * m
		recordMap["zuanShi"] = 999999999 * m
		recordMap["tiLi"] = 999 * m
	case 7:
		// 套餐7: 武器、副武器、关卡等级清至1级, 8元
		recordMap["lDamage"] = 1
		recordMap["lCount"] = 1
		recordMap["levelFuCount"] = "[1,1,1,1,1,1,1,1,1,1]"
		recordMap["levelFuDamage"] = "[1,1,1,1,1,1,1,1,1,1]"
		recordMap["level"] = 1
	// 套餐8: 任意调整关卡等级, 5元
	case 8:
		recordMap["level"] = *level
	default:
		log.Fatal("请选择正确的套餐")
	}

	recordMap["sign"] = SignDataMap(recordMap)
	recordJsonStr, _ := json.Marshal(recordMap)

	fmt.Println(string(recordJsonStr))

	reqMap := make(map[string]interface{})
	reqMap["plat"] = "wx"
	reqMap["record"] = string(recordJsonStr)
	reqMap["time"] = time.Now().UnixNano() / 1e6
	reqMap["openid"] = *openID
	reqMap["wx_appid"] = appId
	reqMap["wx_secret"] = secret
	reqMap["sign"] = SignMap(reqMap)

	delete(reqMap, "wx_appid")
	delete(reqMap, "wx_secret")
	reqData, _ := json.Marshal(reqMap)

	fmt.Println(string(reqData))

	uploadResult := new(RespData)
	if err := PostWxGame("/api/archive/upload", reqData, uploadResult); err != nil {
		log.Fatal(err)
	}
	if uploadResult.Code == 0 {
		fmt.Println("更新用户信息成功")
	} else {
		fmt.Printf("刷新数据失败,%v", uploadResult)
	}
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
