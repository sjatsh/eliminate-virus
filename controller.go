package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var m int64 = 999999999999999999

type Request struct {
	ID       string `form:"id" binding:"required"` // 用户id
	Package  int64  `form:"p" binding:"required"`  // 套餐种类
	Password string `form:"sp" binding:"required"` // 用户密码
	Level    int64  `form:"l"`                     // 关卡等级
	Cancel   string `form:"c"`                     // 取消相应套餐
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func UploadHandler(c *gin.Context) {

	resp := new(Response)
	defer c.JSON(http.StatusOK, resp)

	req := new(Request)
	if err := c.ShouldBind(req); err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		return
	}

	uploadResult, err := ModifyUser(req.ID, req.Cancel, req.Package, req.Level)
	if err != nil {
		resp.Code = -1
		resp.Msg = err.Error()
		return
	}
	resp.Code = uploadResult.Code

	if uploadResult.Code == 0 {
		resp.Data = "更新用户信息成功"
	} else {
		resp.Msg = fmt.Sprintf("刷新数据失败,%v", uploadResult)
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
