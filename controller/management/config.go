/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package management

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi/fishpond/dao/mysql"
	"github.com/wangyi/fishpond/model"
	"github.com/wangyi/fishpond/util"
	"strconv"
)

/**

 */
func SetConfig(c *gin.Context) {
	action := c.PostForm("action")
	if action == "GET" {
		config := model.Config{}
		err := mysql.DB.Where("id=?", 1).First(&config).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "获取失败")
			return
		}
		util.JsonWrite(c, 200, config, "获取成功")
		return
	}

	if action == "UPDATE" {

		config := model.Config{}
		if BAddress, isExist := c.GetPostForm("b_address"); isExist == true {
			config.BAddress = BAddress
		}

		if CAddress, isExist := c.GetPostForm("c_address"); isExist == true {
			config.CAddress = CAddress
		}

		if BKey, isExist := c.GetPostForm("b_key"); isExist == true {
			config.BKey = BKey
		}

		if BMnemonic, isExist := c.GetPostForm("b_mnemonic"); isExist == true {
			config.BMnemonic = BMnemonic
		}

		if RevenueModel, isExist := c.GetPostForm("revenue_model"); isExist == true {
			config.RevenueModel, _ = strconv.Atoi(RevenueModel)
		}

		if IfNeedInCodeString, isExist := c.GetPostForm("IfNeedInCode"); isExist == true {
			config.IfNeedInCode, _ = strconv.Atoi(IfNeedInCodeString)
		}


		if AddMoneyMode, isExist := c.GetPostForm("add_money_mode"); isExist == true {
			config.AddMoneyMode, _ = strconv.Atoi(AddMoneyMode)
		}

		err := mysql.DB.Model(&model.Config{}).Where("id=1").Update(&config).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "修改失败")
			return
		}

		util.JsonWrite(c, 200, nil, "修改成功")
		return

	}

}
