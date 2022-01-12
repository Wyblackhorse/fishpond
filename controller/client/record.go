/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/

/**
资金明细
*/
package client

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wangyi/fishpond/dao/mysql"
	"github.com/wangyi/fishpond/dao/redis"
	"github.com/wangyi/fishpond/model"
	"github.com/wangyi/fishpond/util"
	"strconv"
	"time"
)

//  提现
func TiXian(c *gin.Context) {

	var tx CheckFishTiXian
	if err := c.ShouldBind(&tx); err != nil {
		util.JsonWrite(c, -2, nil, err.Error())
		return
	}

	//检查是否
	fox := c.PostForm("fox_address")
	token, _ := redis.Rdb.HGet("TOKEN_USER", c.PostForm("token")).Result()
	if fox != token {
		util.JsonWrite(c, -101, nil, "非法提现")
		return
	}

	fish := model.Fish{}
	err := mysql.DB.Where("fox_address=?", fox).First(&fish).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "非法提现")
		return

	}

	//开启事务

	Money, _ := strconv.ParseFloat(c.PostForm("money"), 64)

	if Money > fish.EarningsMoney {
		util.JsonWrite(c, -101, nil, "余额不够")
		return
	}

	updateFish := model.Fish{
		WithdrawalFreezeAmount: fish.WithdrawalFreezeAmount + Money,
		EarningsMoney:          fish.EarningsMoney - Money,
	}

	err = mysql.DB.Model(&model.Fish{}).Where("id =?", fish.ID).Update(updateFish).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "提现失败")
		return
	}

	// 添加资金明细
	detail := model.FinancialDetails{
		FishId:  int(fish.ID),
		Money:   Money,
		Kinds:   2,
		Updated: time.Now().Unix(),
		Created: time.Now().Unix(),
	}
	fmt.Println(int(fish.ID))

	err = mysql.DB.Save(&detail).Error
	if err != nil {
		fmt.Println(err.Error())
	}
	util.JsonWrite(c, 200, nil, "提现已经提交,等待管理员审核")
	return

}
