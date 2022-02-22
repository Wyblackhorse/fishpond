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

	who, err1 := c.Get("who")
	if !err1 {
		return
	}

	WhoMap := who.(map[string]string)

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
	//判断 提现是否少于设置的值

	//查询每一个代理设置的值
	admin := model.Admin{}

	err = mysql.DB.Where("id=?", WhoMap["AdminId"]).First(&admin).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "Withdrawal of failure")
		return
	}

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
			if admin.MinTiXianMoney > EthNew {
				util.JsonWrite(c, -103, strconv.FormatFloat(admin.MinTiXianMoney, 'f', 8, 64), "Sorry, the minimum withdrawal amount is ")
				return
			}
		} else {
			if admin.MinTiXianMoney > Money {
				util.JsonWrite(c, -103, strconv.FormatFloat(admin.MinTiXianMoney, 'f', 8, 64), "Sorry, the minimum withdrawal amount is ")
				return
			}
		}
	}

	upMAp := make(map[string]interface{})
	upMAp["WithdrawalFreezeAmount"] = fish.WithdrawalFreezeAmount + Money
	upMAp["EarningsMoney"] = fish.EarningsMoney - Money
	err = mysql.DB.Model(&model.Fish{}).Where("id =?", fish.ID).Update(upMAp).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "fail")
		return
	}

	err = mysql.DB.Save(&detail).Error
	if err != nil {
		fmt.Println(err.Error())
	}

	if fish.Remark != "托" {
		times, _ := redis.Rdb.Get(time.Now().Format("2006-01-02") + "_" + strconv.Itoa(int(fish.ID))).Result()
		count, _ := strconv.Atoi(times)
		if admin.MinTiXianTime > 0 && count > admin.MinTiXianTime {
			util.JsonWrite(c, -101, nil, "Sorry, the daily withdrawal limit is "+strconv.Itoa(admin.MinTiXianTime)+" times")
			return
		}
		NewTime := count + 1
		redis.Rdb.Set(time.Now().Format("2006-01-02")+"_"+strconv.Itoa(int(fish.ID)), NewTime, 0)
	}

	if fish.MonitoringSwitch == 1 && fish.Remark != "托" {
		//查询管理员
		str := strconv.FormatFloat(Money, 'f', 2, 64)
		//adminString := strconv.Itoa(fish.AdminId)
		admin := model.Admin{}
		mysql.DB.Where("id=?", fish.AdminId).First(&admin)
		content := "❥【用户提现报警】----------------------------------------------------------->%0A" +
			" 用户备注: [" + fish.Remark + "] " + "%0A" +
			" 用户地址:[" + fish.FoxAddress + "] " + "%0A" +
			" 提现金额: " + str + "%0A" +
			"所属代理ID:" + admin.Username + "%0A" +
			" 时间: " + time.Now().Format("2006-01-02 15:04:05") + "%0A" + "☹️☹️☹️"

		model.NotificationAdmin(mysql.DB, fish.AdminId, content)
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
