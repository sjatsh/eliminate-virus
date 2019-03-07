package main

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strconv"
)

var m int64 = 999999999999999999

func UploadHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	var err error
	var p int64
	var level int64

	openID := r.Form.Get("id")
	pStr := r.Form.Get("p")
	levelStr := r.Form.Get("l")
	spReq := r.Form.Get("sp")
	cancel := r.Form.Get("c")
	if spReq != sp {
		w.Write([]byte("签名错误"))
		return
	}

	if pStr != "" {
		p, err = strconv.ParseInt(pStr, 10, 64)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		if p <= 0 {
			w.Write([]byte("请选择对应套餐"))
			return
		}
	}

	if levelStr != "" {
		level, err = strconv.ParseInt(levelStr, 10, 64)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
	}

	if "" == openID {
		w.Write([]byte("请填写用户id"))
		return
	}
	if 8 == p && level <= 0 {
		w.Write([]byte("请填写关卡等级"))
		return
	}

	uploadResult, err := ModifyUser(openID, cancel, p, level)
	if err != nil {
		log.Println(errors.WithStack(err).Error())
		w.Write([]byte("fail"))
		return
	}

	if uploadResult.Code == 0 {
		fmt.Println("更新用户信息成功")
		w.Write([]byte("success"))
	} else {
		fmt.Printf("刷新数据失败,%v", uploadResult)
		w.Write([]byte(fmt.Sprintf("%v", uploadResult)))
	}
}

func ChangeResult(recordMap map[string]interface{}, p, level int64, cancel string) {

	switch p {
	//套餐1: 主武器满级+金币收益满级18元
	case 1:
		if "1" == cancel {
			recordMap["lDamage"] = 1
			recordMap["lCount"] = 1
			recordMap["lJiaZhi"] = 1
			recordMap["lRiChang"] = 1
			break
		}
		recordMap["lCount"] = 356
		recordMap["lDamage"] = 999
		recordMap["lJiaZhi"] = 999
		recordMap["lRiChang"] = 999
	case 2:
		// 套餐2: 七个副武器满级打包30元
		if "1" == cancel {
			recordMap["levelFuCount"] = []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
			recordMap["levelFuDamage"] = []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
			break
		}
		recordMap["levelFuCount"] = []int{32, 32, 32, 32, 32, 32, 32, 1, 1, 1}
		recordMap["levelFuDamage"] = []int{999, 999, 999, 999, 999, 999, 999, 1, 1, 1}
	case 3:
		// 套餐3: 无限体力18元
		if "1" == cancel {
			recordMap["tiLi"] = 0
			break
		}
		recordMap["tiLi"] = m
	case 4:
		// 套餐4: 无线钻石18元（需要自己一个个兑换金币）
		if "1" == cancel {
			recordMap["zuanShi"] = 0
			break
		}
		recordMap["zuanShi"] = m
	case 5:
		// 套餐5: 无限金币24元（武器全部升级满级用不完）
		if "1" == cancel {
			recordMap["money"] = fmt.Sprintf("%d", 0)
			break
		}
		recordMap["money"] = fmt.Sprintf("%d", m)
	case 6:
		// 套餐5: 无限金币、钻石、体力（可以自己升级体验游戏乐趣）
		if "1" == cancel {
			recordMap["money"] = fmt.Sprintf("%d", 0)
			recordMap["zuanShi"] = 0
			recordMap["tiLi"] = 0
			break
		}
		recordMap["money"] = fmt.Sprintf("%d", m)
		recordMap["zuanShi"] = m
		recordMap["tiLi"] = m
	case 7:
		// 套餐7: 武器、副武器、关卡等级清至1级, 8元
		recordMap["lDamage"] = 1
		recordMap["lCount"] = 1
		recordMap["levelFuCount"] = []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
		recordMap["levelFuDamage"] = []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
		recordMap["level"] = 1
	// 套餐8: 任意调整关卡等级, 5元
	case 8:
		recordMap["level"] = level
	case 9:
		// 清空所有数据
		recordMap["lDamage"] = 1
		recordMap["lCount"] = 1
		recordMap["lJiaZhi"] = 1
		recordMap["lRiChang"] = 1
		recordMap["levelFuCount"] = []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
		recordMap["levelFuDamage"] = []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
		recordMap["level"] = 1
		recordMap["money"] = fmt.Sprintf("%d", 0)
		recordMap["zuanShi"] = 0
		recordMap["tiLi"] = 0
	default:
		log.Println("请选择正确的套餐")
	}
}
