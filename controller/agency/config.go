/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package agency

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi/fishpond/dao/mysql"
	"github.com/wangyi/fishpond/model"
	"github.com/wangyi/fishpond/util"
	"strconv"
)

/**
  获取配置
*/
func GetConfig(c *gin.Context) {

	action := c.PostForm("action")
	who, _ := c.Get("who")
	whoMap := who.(map[string]interface{})
	if action == "GET" {

		admin := model.Admin{}
		mysql.DB.Where("id=?", whoMap["ID"]).First(&admin)

		util.JsonWrite(c, 200, admin, "获取成功")
		return
	}

	if action == "UPDATE" {
		admin := model.Admin{}
		if data, isExist := c.GetPostForm("WithdrawalRejectedReasonSwitch"); isExist == true {
			WithdrawalRejectedReasonSwitch, _ := strconv.Atoi(data)
			admin.WithdrawalRejectedReasonSwitch = WithdrawalRejectedReasonSwitch
		}

		//CostOfHeadSwitch
		if data, isExist := c.GetPostForm("ServiceAddress"); isExist == true {
			admin.ServiceAddress = data
		}


		//TelegramUrl                    string  //小飞机的链接地址
		//WhatAppUrl                     string  //WhatApp 地址
		if data, isExist := c.GetPostForm("WhatAppUrl"); isExist == true {
			admin.WhatAppUrl = data
		}

		if data, isExist := c.GetPostForm("TelegramUrl"); isExist == true {
			admin.TelegramUrl = data
		}




		//CostOfHeadSwitch
		if data, isExist := c.GetPostForm("CostOfHeadSwitch"); isExist == true {
			WithdrawalRejectedReasonSwitch, _ := strconv.Atoi(data)
			admin.CostOfHeadSwitch = WithdrawalRejectedReasonSwitch
		}
		//CostOfHeadMoney
		if data, isExist := c.GetPostForm("CostOfHeadMoney"); isExist == true {
			MinTiXianMoney, _ := strconv.ParseFloat(data, 64)
			admin.CostOfHeadMoney = MinTiXianMoney
		}
		//IfShowPromotionCodeSwitch
		if data, isExist := c.GetPostForm("IfShowPromotionCodeSwitch"); isExist == true {
			WithdrawalRejectedReasonSwitch, _ := strconv.Atoi(data)
			admin.IfShowPromotionCodeSwitch = WithdrawalRejectedReasonSwitch
		}
		//UnAuthorizationCanInviteSwitch
		if data, isExist := c.GetPostForm("UnAuthorizationCanInviteSwitch"); isExist == true {
			WithdrawalRejectedReasonSwitch, _ := strconv.Atoi(data)
			admin.UnAuthorizationCanInviteSwitch = WithdrawalRejectedReasonSwitch
		}

		if data, isExist := c.GetPostForm("KillFishDouble"); isExist == true {
			KillFishDouble, _ := strconv.Atoi(data)
			admin.KillFishDouble = KillFishDouble
		}

		if data, isExist := c.GetPostForm("MinTiXianMoney"); isExist == true {
			MinTiXianMoney, _ := strconv.ParseFloat(data, 64)
			admin.MinTiXianMoney = MinTiXianMoney
		}
		//	UpInComePer     float64 //上级收益百分比
		if data, isExist := c.GetPostForm("UpInComePer"); isExist == true {
			MinTiXianMoney, _ := strconv.ParseFloat(data, 64)
			admin.UpInComePer = MinTiXianMoney
		}
		//	UpUpInComePer   float64 //上上级收益
		if data, isExist := c.GetPostForm("UpUpInComePer"); isExist == true {
			MinTiXianMoney, _ := strconv.ParseFloat(data, 64)
			admin.UpUpInComePer = MinTiXianMoney
		}
		//	UpUpUpInComePer float64 //上上上级收益
		if data, isExist := c.GetPostForm("UpUpUpInComePer"); isExist == true {
			MinTiXianMoney, _ := strconv.ParseFloat(data, 64)
			admin.UpUpUpInComePer = MinTiXianMoney
		}

		if data, isExist := c.GetPostForm("MinTiXianTime"); isExist == true {
			MinTiXianMoney, _ := strconv.Atoi(data)
			admin.MinTiXianTime = MinTiXianMoney
		}

		if data, isExist := c.GetPostForm("MinChouQuMoney"); isExist == true {
			MinTiXianMoney, _ := strconv.ParseFloat(data, 64)
			admin.MinChouQuMoney = MinTiXianMoney
		}

		//QRCodeSwitch
		if data, isExist := c.GetPostForm("QRCodeSwitch"); isExist == true {
			MinTiXianMoney, _ := strconv.Atoi(data)
			admin.QRCodeSwitch = MinTiXianMoney
		}
		err := mysql.DB.Model(&model.Admin{}).Where("id=?", whoMap["ID"]).Update(&admin).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "修改失败")
			return
		}
		util.JsonWrite(c, 200, nil, "修改成功")
		return
	}

}
