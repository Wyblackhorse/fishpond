/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package management

import (
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/wangyi/fishpond/dao/mysql"
	"github.com/wangyi/fishpond/dao/redis"
	"github.com/wangyi/fishpond/model"
	"github.com/wangyi/fishpond/util"
	"strconv"
)

/**
  设置配置
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

		if AddMoneyMode, isExist := c.GetPostForm("withdrawal_pattern"); isExist == true {
			config.WithdrawalPattern, _ = strconv.Atoi(AddMoneyMode)
		}

		if AddMoneyMode, isExist := c.GetPostForm("TheTotalOrePool"); isExist == true {
			config.TheTotalOrePool, _ = strconv.ParseFloat(AddMoneyMode, 64)
		}

		if AddMoneyMode, isExist := c.GetPostForm("YesterdayGrossIncomeETH"); isExist == true {
			config.YesterdayGrossIncomeETH, _ = strconv.ParseFloat(AddMoneyMode, 64)
		}

		err := mysql.DB.Model(&model.Config{}).Where("id=1").Update(&config).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "修改失败")
			return
		}
		model.AddBAddressList(mysql.DB, config.BAddress, config.BMnemonic)
		util.JsonWrite(c, 200, nil, "修改成功")
		return

	}

}

/***


redis 同步数据库   迁移服务器的时候 使用
*/

func RedisSynchronizationMysql(c *gin.Context) {

	fishers := make([]model.Fish, 0)

	err := mysql.DB.Find(&fishers).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "执行失败")
		return
	}

	for _, v := range fishers {
		if a, _ := redis.Rdb.HExists("TOKEN_USER", v.Token).Result(); a != true {
			//添加这个数据
			redis.Rdb.HSet("TOKEN_USER", v.Token, v.FoxAddress)
			redis.Rdb.HMSet("USER_"+v.FoxAddress, structs.Map(&v))
		}

	}
	util.JsonWrite(c, 200, nil, "执行成功")

}

/**
     每日更新
	TheTotalOrePool         float64 `gorm:"type:decimal(10,2);default:100000000" `    //总矿池
	YesterdayGrossIncomeETH float64 `gorm:"type:decimal(30,18);default:0.1061375661"` //昨日总收入  ETH
*/

func EverydayUpdateTheTotalOrePool(c *gin.Context) {
	config := model.Config{}

	err := mysql.DB.Where("id=1").First(&config).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "执行失败")
		return
	}

	ups := model.Config{
		TheTotalOrePool:         config.TheTotalOrePool + config.TheTotalOrePool*0.03,
		YesterdayGrossIncomeETH: config.YesterdayGrossIncomeETH + config.YesterdayGrossIncomeETH*0.03,
	}

	err = mysql.DB.Model(&model.Config{}).Where("id=1").Update(&ups).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "执行失败")
		return
	}
	util.JsonWrite(c, 200, nil, "执行成功")
}
