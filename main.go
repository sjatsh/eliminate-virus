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

func main() {

	coin := flag.String("coin", "99999999999", "金币数")
	zuanShi := flag.Int("zuanshi", 99999999999, "钻石数量")
	openID := flag.String("open_id", "", "open_id")
	isSoundOff := flag.Bool("sound_off", false, "true:关闭音乐,false:开启音乐")
	level := flag.Int("level", 200, "关卡等级")
	lDamage := flag.Int("damage", 999, "主武器火力等级")
	lCount := flag.Int("count", 356, "主武器射速")
	levelFuCountStr := flag.String("levelFuCount", "[32,32,32,32,32,32,32,1,1,1]", "副武器强度，更改强度需要先删除本地小程序")
	levelFuDamageStr := flag.String("levelFuDamage", "[999,999,999,999,999,999,999,1,1,1]", "副武器火力")
	jiazhi := flag.Int("jiazhi", 999, "金币价值等级")
	richang := flag.Int("richang", 999, "日常收益等级")
	tili := flag.Int("tili", 100, "体力值")
	playCount := flag.Int("playCount", 1, "游戏次数")
	shareCount := flag.Int("shareCount", 1, "分享次数")
	videoCount := flag.Int("videoCount", 1, "广告观看次数")
	flag.Parse()

	levelFuCount := make([]int, 0)
	if err := json.Unmarshal([]byte(*levelFuCountStr), &levelFuCount); err != nil {
		log.Fatalf("副武器强度配置有误:%s", err)
	}
	levelFuDamage := make([]int, 0)
	if err := json.Unmarshal([]byte(*levelFuDamageStr), &levelFuDamage); err != nil {
		log.Fatalf("副武器火力配置有误:%s", err)
	}

	recordMap := make(map[string]interface{})
	recordMap["uid"] = *openID
	recordMap["isSoundOff"] = *isSoundOff
	recordMap["level"] = *level
	recordMap["lDamage"] = *lDamage
	recordMap["lCount"] = *lCount
	recordMap["lJiaZhi"] = *jiazhi
	recordMap["lRiChang"] = *richang
	recordMap["curFu"] = 2
	recordMap["levelFuCount"] = levelFuCount
	recordMap["levelFuDamage"] = levelFuDamage
	recordMap["getTime"] = "26,12,23,51"
	recordMap["bgIndex"] = 1
	recordMap["money"] = *coin
	recordMap["tipFU"] = false
	recordMap["isGuide"] = false
	recordMap["tiLi"] = *tili
	recordMap["tiLiBackTime"] = 65365793
	recordMap["today"] = time.Now().Day()
	recordMap["playCount"] = *playCount
	recordMap["shareCount"] = *shareCount
	recordMap["videoCount"] = *videoCount
	recordMap["isGuanZhu"] = true
	recordMap["isShouCang"] = true
	recordMap["tryFuCount"] = 0
	recordMap["pos"] = ""
	recordMap["posUpdate"] = 0
	recordMap["zuanShi"] = *zuanShi

	// 取出所有的键，并排序
	keys := make([]string, 0, len(recordMap))
	for key := range recordMap {
		if key == "sign" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// 拼接参数
	ignS := ""
	for i, j := 0, len(keys); i < j; i++ {
		ignS += fmt.Sprintf("%v", recordMap[keys[i]])
	}
	bs := md5.Sum(bytes.NewBufferString(ignS).Bytes())
	ms := strings.ToLower(hex.EncodeToString(bs[:]))
	recordMap["sign"] = ms

	recordJsonStr, _ := json.Marshal(recordMap)

	reqMap := make(map[string]interface{})
	reqMap["plat"] = "wx"
	reqMap["record"] = string(recordJsonStr)
	reqMap["time"] = time.Now().UnixNano() / 1e6
	reqMap["openid"] = *openID
	reqMap["wx_appid"] = appId
	reqMap["wx_secret"] = secret

	reqkeys := make([]string, 0, len(reqMap))
	for key := range reqMap {
		if key == "sign" {
			continue
		}
		reqkeys = append(reqkeys, key)
	}
	sort.Strings(reqkeys)

	reqSignEle := make([]string, 0)
	for i, j := 0, len(reqkeys); i < j; i++ {
		reqSignEle = append(reqSignEle, reqkeys[i]+"="+fmt.Sprintf("%v", reqMap[reqkeys[i]]))
	}
	reqbs := md5.Sum(bytes.NewBufferString(strings.Join(reqSignEle, "&")).Bytes())
	reqms := strings.ToLower(hex.EncodeToString(reqbs[:]))
	reqMap["sign"] = reqms
	delete(reqMap, "wx_appid")
	delete(reqMap, "wx_secret")

	reqData, _ := json.Marshal(reqMap)

	fmt.Printf("req string: %s\n", string(reqData))
	resp, err := http.Post("https://wxwyjh.chiji-h5.com/api/archive/upload", "application/json", bytes.NewReader(reqData))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	result := struct {
		Data map[string]interface{} `json:"data"`
		Code int                    `json:"code"`
	}{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatal(err)
	}
	if result.Code == 0 {
		fmt.Println("刷新成功")
	} else {
		fmt.Printf("刷新数据失败,%v", result)
	}

}
