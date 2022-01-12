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

func GetVipEarnings(c *gin.Context)  {

	_, err1 := c.Get("who")
	if !err1 {
		return
	}
	vip := make([]model.VipEarnings, 0)
	err := mysql.DB.Find(&vip).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "获取失败")
		return
	}
	util.JsonWrite(c, 200, vip, "获取成功")
	return

}
