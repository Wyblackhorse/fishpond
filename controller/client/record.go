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
	"net/http"
	"strconv"
	"time"
)

//  提现
func TiXian(c *gin.Context) {

	_, err1 := c.Get("who")
	if !err1 {
		return
	}

	var tx CheckFishTiXian
	if err := c.ShouldBind(&tx); err != nil {
		util.JsonWrite(c, -2, nil, err.Error())
		return
	}

	//检查是否
	fox := c.PostForm("fox_address")
	token, _ := redis.Rdb.HGet("TOKEN_USER", c.PostForm("token")).Result()
	if fox != token {
		util.JsonWrite(c, -101, nil, "Withdrawal of failure")
		return
	}

	fish := model.Fish{}
	err := mysql.DB.Where("fox_address=?", fox).First(&fish).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "Withdrawal of failure")
		return
	}

	//开启事务
	Money, _ := strconv.ParseFloat(c.PostForm("money"), 64) //提现
	if Money > fish.EarningsMoney {
		util.JsonWrite(c, -101, nil, "The balance is not enough")
		return
	}
	//换算成ETH

	// 添加资金明细
	detail := model.FinancialDetails{
		FishId:  int(fish.ID),
		Money:   Money,
		Kinds:   2,
		Updated: time.Now().Unix(),
		Created: time.Now().Unix(),
	}
	if kinds, isExist := c.GetPostForm("kinds"); isExist == true {
		if kinds == "2" { //提现 ETH
			//获取eth 汇率
			HuiLV, _ := redis.Rdb.Get("ETHTOUSDT").Result()
			HH, _ := strconv.ParseFloat(HuiLV, 64) //提现
			EthNew := Money / HH
			detail.MoneyEth = EthNew
			detail.Pattern = 2
			detail.TheExchangeRateAtThatTime = HH
		}
	}

	updateFish := model.Fish{
		WithdrawalFreezeAmount: fish.WithdrawalFreezeAmount + Money,
		EarningsMoney:          fish.EarningsMoney - Money,
	}
	err = mysql.DB.Model(&model.Fish{}).Where("id =?", fish.ID).Update(updateFish).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "fail")
		return
	}

	err = mysql.DB.Save(&detail).Error
	if err != nil {
		fmt.Println(err.Error())
	}
	util.JsonWrite(c, 200, nil, "The request has been submitted pending background review")
	return
}

/**

  获取自己 收益订单
*/

func GetEarnings(c *gin.Context) {

	data, err1 := c.Get("who")
	if !err1 {
		return
	}
	a := data.(map[string]string)
	page, _ := strconv.Atoi(c.PostForm("page"))
	limit, _ := strconv.Atoi(c.PostForm("limit"))
	var total int = 0
	Db := mysql.DB
	recodes := make([]model.FinancialDetails, 0)
	if status, isExist := c.GetPostForm("kinds"); isExist == true {
		status, _ := strconv.Atoi(status)
		Db = Db.Where("kinds=?", status)
	}
	Db.Table("financial_details").Count(&total)
	Db = Db.Model(&recodes).Where("fish_id=?", a["ID"]).Offset((page - 1) * limit).Limit(limit).Order("created desc")
	if err := Db.Find(&recodes).Error; err != nil {
		util.JsonWrite(c, -101, nil, err.Error())
		return
	}

	//获取最新的汇率
	exchange, errR := redis.Rdb.Get("ETHTOUSDT").Result()

	if errR != nil {
		util.JsonWrite(c, -101, nil, "Exchange rate acquisition failure")
		return

	}

	newExchange, _ := strconv.ParseFloat(exchange, 64)

	for k, _ := range recodes {
		recodes[k].ETH = recodes[k].Money / newExchange
	}

	c.JSON(http.StatusOK, gin.H{
		"code":   1,
		"count":  total,
		"result": recodes,
	})
	return
}
