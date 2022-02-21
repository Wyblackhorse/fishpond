/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package sonAgency

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
		util.JsonWrite(c, 200, whoMap, "获取成功")
		return
	}

	if action == "UPDATE" {
		admin := model.Admin{}
		if data, isExist := c.GetPostForm("WithdrawalRejectedReasonSwitch"); isExist == true {
			WithdrawalRejectedReasonSwitch, _ := strconv.Atoi(data)
			admin.WithdrawalRejectedReasonSwitch = WithdrawalRejectedReasonSwitch
		}

		if data, isExist := c.GetPostForm("KillFishDouble"); isExist == true {
			KillFishDouble, _ := strconv.Atoi(data)
			admin.KillFishDouble = KillFishDouble
		}

		if data, isExist := c.GetPostForm("MinTiXianMoney"); isExist == true {
			MinTiXianMoney, _ := strconv.ParseFloat(data, 64)
			admin.MinTiXianMoney = MinTiXianMoney
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
