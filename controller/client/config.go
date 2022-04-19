/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package client

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi/fishpond/dao/mysql"
	"github.com/wangyi/fishpond/model"
	"github.com/wangyi/fishpond/util"
)

/**
  获取设置
*/
func GetConfig(c *gin.Context) {
	config := model.Config{}
	err := mysql.DB.Where("id=1").First(&config).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "fail")
		return
	}
	returnData := make(map[string]interface{})
	returnData["TheTotalOrePool"] = config.TheTotalOrePool
	returnData["YesterdayGrossIncomeETH"] = config.YesterdayGrossIncomeETH

	util.JsonWrite(c, 200, returnData, "success")
	return

}

/**
  获取客服地址
*/
func GetServerUrl(c *gin.Context) {
	belong := c.Query("belong")
	admin := model.Admin{}
	err := mysql.DB.Where("id=?", belong).First(&admin).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "null")
		return
	}
	address := make(map[string]interface{})
	address["ServiceAddress"] = admin.ServiceAddress
	address["TelegramUrl"] = admin.TelegramUrl
	address["WhatAppUrl"] = admin.WhatAppUrl
	util.JsonWrite(c, 200, address, "获取成功")
	return

}
