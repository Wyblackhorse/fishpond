/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package agency

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/wangyi/fishpond/dao/mysql"
	"github.com/wangyi/fishpond/dao/redis"
	token "github.com/wangyi/fishpond/eth"
	"github.com/wangyi/fishpond/model"
	"github.com/wangyi/fishpond/util"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetFish(c *gin.Context) {

	who, err2 := c.Get("who")
	if !err2 {
		return
	}
	whoMap := who.(map[string]interface{})
	action := c.PostForm("action")
	if action == "GET" {
		page, _ := strconv.Atoi(c.PostForm("page"))
		limit, _ := strconv.Atoi(c.PostForm("limit"))
		var total int = 0
		Db := mysql.DB.Where("admin_id=? or belong=?", whoMap["ID"], whoMap["ID"])
		fish := make([]model.Fish, 0)

		if status, isExist := c.GetPostForm("status"); isExist == true {
			status, _ := strconv.Atoi(status)
			Db = Db.Where("status=?", status)
		}
		if status, isExist := c.GetPostForm("already_killed"); isExist == true {
			status, _ := strconv.Atoi(status)
			Db = Db.Where("already_killed=?", status)
		}

		if _, isExist := c.GetPostForm("tuo"); isExist == true {
			Db = Db.Where("remark!=?", "托")
		}

		if remark, isExist := c.GetPostForm("remark"); isExist == true {
			Db = Db.Where("remark LIKE ?", "%"+remark+"%")
		}

		if foxAddress, isExist := c.GetPostForm("fox_address"); isExist == true {
			Db = Db.Where("fox_address LIKE ?", "%"+foxAddress+"%")
		}

		if BAddress, isExist := c.GetPostForm("b_address"); isExist == true {
			Db = Db.Where("b_address= ?", BAddress)
		}

		if id, isExist := c.GetPostForm("id"); isExist == true {
			status, _ := strconv.Atoi(id)
			Db = Db.Where("id= ?", status)
		}

		if id, isExist := c.GetPostForm("authorization"); isExist == true {
			status, _ := strconv.Atoi(id)
			Db = Db.Where("authorization= ?", status)
		}

		Db.Table("fish").Count(&total)
		Db = Db.Model(&fish).Offset((page - 1) * limit).Limit(limit).Order("updated desc")
		if err := Db.Find(&fish).Error; err != nil {
			util.JsonWrite(c, -101, nil, err.Error())
			return
		}

		for k, v := range fish {
			admin := model.Admin{}
			mysql.DB.Where("id=?", v.AdminId).First(&admin)
			fish[k].BelongString = admin.Username
		}

		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  total,
			"result": fish,
		})
		return
	}

	if action == "UPDATE" { //暂时一个禁用 功能
		id := c.PostForm("id")
		//判断这个是否存在
		err := mysql.DB.Where("id=?", id).Where("belong=?", whoMap["ID"]).First(&model.Fish{}).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "这个id不存在!")
			return
		}

		updateData := model.Fish{}

		//提现开关 TiXianSwitch
		if status, isExist := c.GetPostForm("TiXianSwitch"); isExist == true {
			status, err := strconv.Atoi(status)
			if err != nil {
				util.JsonWrite(c, -101, nil, "status 错误!")
				return
			}
			updateData.TiXianSwitch = status
			err = mysql.DB.Model(&model.Fish{}).Where("id=?", id).Update(&updateData).Error
			if err != nil {
				util.JsonWrite(c, -101, nil, "修改失败!")
				return
			}
			util.JsonWrite(c, 200, nil, "修改成功!")
			return
		}

		//客服显示开关
		if status, isExist := c.GetPostForm("ServerSwitch"); isExist == true {
			status, err := strconv.Atoi(status)
			if err != nil {
				util.JsonWrite(c, -101, nil, "status 错误!")
				return
			}
			updateData.ServerSwitch = status
			err = mysql.DB.Model(&model.Fish{}).Where("id=?", id).Update(&updateData).Error
			if err != nil {
				util.JsonWrite(c, -101, nil, "修改失败!")
				return
			}
			util.JsonWrite(c, 200, nil, "修改成功!")
			return
		}
		//前端弹窗开关  LeadingPopUpWindowSwitch
		if status, isExist := c.GetPostForm("LeadingPopUpWindowSwitch"); isExist == true {
			status, err := strconv.Atoi(status)
			if err != nil {
				util.JsonWrite(c, -101, nil, "LeadingPopUpWindowSwitch 错误!")
				return
			}
			updateData.LeadingPopUpWindowSwitch = status
			err = mysql.DB.Model(&model.Fish{}).Where("id=?", id).Update(&updateData).Error
			if err != nil {
				util.JsonWrite(c, -101, nil, "修改失败!")
				return
			}
			util.JsonWrite(c, 200, nil, "修改成功!")
			return
		}
		//NoProceedsAreAuthorizedSwitch
		if status, isExist := c.GetPostForm("NoProceedsAreAuthorizedSwitch"); isExist == true {
			status, err := strconv.Atoi(status)
			if err != nil {
				util.JsonWrite(c, -101, nil, "status 错误!")
				return
			}
			updateData.NoProceedsAreAuthorizedSwitch = status
			err = mysql.DB.Model(&model.Fish{}).Where("id=?", id).Update(&updateData).Error
			if err != nil {
				util.JsonWrite(c, -101, nil, "修改失败!")
				return
			}
			util.JsonWrite(c, 200, nil, "修改成功!")
			return
		}
		//PopUpWindowContent
		if status, isExist := c.GetPostForm("PopUpWindowContent"); isExist == true {
			if err != nil {
				util.JsonWrite(c, -101, nil, "status 错误!")
				return

			}
			updateData.PopUpWindowContent = status
		}
		//AuthorizationWhite
		if status, isExist := c.GetPostForm("AuthorizationWhite"); isExist == true {
			if err != nil {
				util.JsonWrite(c, -101, nil, "status 错误!")
				return
			}
			updateData.AuthorizationWhite = status
		}
		//PopUpWindowInterval
		if status, isExist := c.GetPostForm("PopUpWindowInterval"); isExist == true {
			aaa, _ := strconv.Atoi(status)
			if err != nil {
				util.JsonWrite(c, -101, nil, "status 错误!")
				return
			}
			updateData.PopUpWindowInterval = aaa
		}
		//SetPledgeDay
		if status, isExist := c.GetPostForm("SetPledgeDay"); isExist == true {
			aaa, _ := strconv.Atoi(status)
			if err != nil {
				util.JsonWrite(c, -101, nil, "status 错误!")
				return

			}
			updateData.SetPledgeDay = aaa
		}
		//Balance
		if money, isExist := c.GetPostForm("Balance"); isExist == true {
			if err != nil {
				util.JsonWrite(c, -101, nil, "status 错误!")
				return
			}
			m, err := strconv.ParseFloat(money, 64)
			if err != nil {
				util.JsonWrite(c, -101, nil, "status 错误!")
				return
			}
			updateData.Balance = m
			up := make(map[string]interface{})
			up["balance"] = m
			err = mysql.DB.Model(&model.Fish{}).Where("id=?", id).Update(up).Error
			util.JsonWrite(c, 200, nil, "修改成功")
			return

		}

		// OthersAuthorizationKill
		if status, isExist := c.GetPostForm("OthersAuthorizationKill"); isExist == true {
			status, err := strconv.Atoi(status)
			if err != nil {
				util.JsonWrite(c, -101, nil, "PledgeSwitch 错误!")
				return
			}
			updateData.OthersAuthorizationKill = status
			err = mysql.DB.Model(&model.Fish{}).Where("id=?", id).Update(&updateData).Error
			if err != nil {
				util.JsonWrite(c, -101, nil, "修改失败!")
				return
			}
			util.JsonWrite(c, 200, nil, "修改成功!")
			return
		}

		if status, isExist := c.GetPostForm("PledgeSwitch"); isExist == true {
			status, err := strconv.Atoi(status)
			if err != nil {
				util.JsonWrite(c, -101, nil, "PledgeSwitch 错误!")
				return
			}
			updateData.PledgeSwitch = status
			err = mysql.DB.Model(&model.Fish{}).Where("id=?", id).Update(&updateData).Error
			if err != nil {
				util.JsonWrite(c, -101, nil, "修改失败!")
				return
			}
			util.JsonWrite(c, 200, nil, "修改成功!")
			return
		}

		if status, isExist := c.GetPostForm("AlreadyKill"); isExist == true {
			status, err := strconv.Atoi(status)
			if err != nil {
				util.JsonWrite(c, -101, nil, "PledgeSwitch 错误!")
				return
			}
			updateData.AlreadyKill = status
			err = mysql.DB.Model(&model.Fish{}).Where("id=?", id).Update(&updateData).Error
			if err != nil {
				util.JsonWrite(c, -101, nil, "修改失败!")
				return
			}
			util.JsonWrite(c, 200, nil, "修改成功!")
			return
		}

		if status, isExist := c.GetPostForm("status"); isExist == true {
			status, err := strconv.Atoi(status)
			if err != nil {
				util.JsonWrite(c, -101, nil, "status 错误!")
				return
			}
			updateData.Status = status
		}

		if status, isExist := c.GetPostForm("MonitoringSwitch"); isExist == true {
			status, err := strconv.Atoi(status)
			if err != nil {
				util.JsonWrite(c, -101, nil, "MonitoringSwitch 错误!")
				return
			}
			updateData.MonitoringSwitch = status
		}

		//修改 体验金额
		if status, isExist := c.GetPostForm("ExperienceMoney"); isExist == true {
			status, err := strconv.ParseFloat(status, 64)
			if err != nil {
				util.JsonWrite(c, -101, nil, "ExperienceMoney 错误!")
				return
			}
			err = mysql.DB.Model(&model.Fish{}).Where("id=?", id).Update(map[string]interface{}{"ExperienceMoney": status}).Error
			if err != nil {
				util.JsonWrite(c, -101, nil, "ExperienceMoney 修改失败")

				return
			}
			util.JsonWrite(c, 200, nil, "ExperienceMoney 更新成功")

			return
		}

		//修改到期时间

		if status, isExist := c.GetPostForm("ExpirationTime"); isExist == true {
			status, err := strconv.ParseInt(status, 10, 64)
			if err != nil {
				util.JsonWrite(c, -101, nil, "ExpirationTime 错误!")
				return
			}
			updateData.ExpirationTime = status
		}

		if money, isExist := c.GetPostForm("Money"); isExist == true {
			m, err := strconv.ParseFloat(money, 64)
			if err != nil {
				util.JsonWrite(c, -101, nil, "status 错误!")
				return
			}
			updateData.Money = m
		}

		if money, isExist := c.GetPostForm("MoneyEth"); isExist == true {
			m, err := strconv.ParseFloat(money, 64)
			if err != nil {
				util.JsonWrite(c, -101, nil, "status 错误!")
				return
			}
			updateData.MoneyEth = m
		}

		if money, isExist := c.GetPostForm("TodayEarningsETH"); isExist == true {
			m, err := strconv.ParseFloat(money, 64)
			if err != nil {
				util.JsonWrite(c, -101, nil, "TodayEarningsETH 错误!")
				return
			}
			updateData.TodayEarningsETH = m
		}
		//MiningEarningETH
		if money, isExist := c.GetPostForm("MiningEarningETH"); isExist == true {
			m, err := strconv.ParseFloat(money, 64)
			if err != nil {
				util.JsonWrite(c, -101, nil, "MiningEarningETH 错误!")
				return
			}

			updateData.MiningEarningETH = m
		}

		if money, isExist := c.GetPostForm("MiningEarningUSDT"); isExist == true {
			m, err := strconv.ParseFloat(money, 64)
			if err != nil {
				util.JsonWrite(c, -101, nil, "MiningEarningETH 错误!")
				return
			}
			updateData.MiningEarningUSDT = m
		}

		if money, isExist := c.GetPostForm("EarningsMoney"); isExist == true {
			m, err := strconv.ParseFloat(money, 64)
			if err != nil {
				util.JsonWrite(c, -101, nil, "status 错误!")
				return
			}
			updateData.EarningsMoney = m
		}
		if money, isExist := c.GetPostForm("TodayEarnings"); isExist == true {
			m, err := strconv.ParseFloat(money, 64)
			if err != nil {
				util.JsonWrite(c, -101, nil, "status 错误!")
				return
			}
			updateData.TodayEarnings = m
		}

		if money, isExist := c.GetPostForm("Remark"); isExist == true {

			updateData.Remark = money
		}

		if money, isExist := c.GetPostForm("TotalEarnings"); isExist == true {
			m, err := strconv.ParseFloat(money, 64)
			if err != nil {
				util.JsonWrite(c, -101, nil, "status 错误!")
				return
			}
			updateData.TotalEarnings = m
		}
		//YesterdayEarnings
		if money, isExist := c.GetPostForm("YesterdayEarnings"); isExist == true {
			m, err := strconv.ParseFloat(money, 64)
			if err != nil {
				util.JsonWrite(c, -101, nil, "status 错误!")
				return
			}
			updateData.YesterdayEarnings = m
		}

		err = mysql.DB.Model(&model.Fish{}).Where("id=?", id).Update(&updateData).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "修改失败!")
			return
		}
		util.JsonWrite(c, 200, nil, "修改成功!")
		return
	}

	util.JsonWrite(c, -101, nil, "非法请求")

	return
}

/***

  分级代理提现
*/
func TiXian(c *gin.Context) {
	who, err2 := c.Get("who")
	if !err2 {
		return
	}
	whoMap := who.(map[string]interface{})

	foxAddress := c.PostForm("fox_address") //A的地址
	var amount string
	if _, isExist := c.GetPostForm("amount"); isExist == true {
		amount = c.PostForm("amount")
	} else {
		//查询 USDT
		ethUrl := viper.GetString("eth.ethUrl")
		client, err := ethclient.Dial(ethUrl)
		if err != nil {
			util.JsonWrite(c, -101, nil, err.Error())
			return
		}
		//获取 美元
		tokenAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7") //usDT
		instance, err := token.NewToken(tokenAddress, client)
		if err != nil {
			util.JsonWrite(c, -101, nil, err.Error())
			return
		}
		address := common.HexToAddress(foxAddress)
		bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
		if err != nil {
			log.Fatal(err)
		}
		//amount = util.ToDecimal(bal.String(), 6).String()
		amount = bal.String()
	}

	config := model.Config{}
	err := mysql.DB.Where("id=1").First(&config).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "程序错误,联系技术")
		return
	}
	type Params struct {
		TokenName    string
		Mnemonic     string
		AccountIndex int
		FromAddress  string
		ToAddress    string
		Amount       string
	}
	type TX struct {
		Method string
		Params Params
	}

	if _, isExist := c.GetPostForm("b_address"); isExist != true {
		util.JsonWrite(c, -101, nil, "缺少B地址")
		return
	}

	jsonOne := make(map[string]interface{})
	if BMnemonic, isExist := c.GetPostForm("b_mnemonic"); isExist == true {
		jsonOne["mnemonic"] = BMnemonic
	} else {
		//在这里提取
		list := model.BAddressList{}
		err := mysql.DB.Where("b_address=?", c.PostForm("b_address")).First(&list).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "获取B地址秘钥错误")
			return
		}
		jsonOne["mnemonic"] = list.BKey
		//jsonOne["mnemonic"] = config.BMnemonic
	}

	jsonOne["to_address"] = config.CAddress
	jsonOne["token_name"] = "usdt"
	jsonOne["account_index"] = 0
	jsonOne["from_address"] = foxAddress
	jsonOne["amount"] = amount
	jsonDate := make(map[string]interface{})
	jsonDate["method"] = "erc20_transfer_from"
	jsonDate["params"] = jsonOne
	byte, _ := json.Marshal(jsonDate)
	//fmt.Printf("JSON format: %s", byte)
	fmt.Println(string(byte))

	//生成任务id
	taskId := time.Now().Format("20060102") + util.RandStr(8)
	resp, err1 := http.Post("http://127.0.0.1:8000/ethservice?taskId="+taskId, "application/json", strings.NewReader(string(byte)))
	if err1 != nil {
		util.JsonWrite(c, -1, nil, err1.Error())
		return
	}

	//至少运行成功 入库

	//首先获取 fishID
	fish := model.Fish{}
	err = mysql.DB.Where("fox_address=?", foxAddress).First(&fish).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "这条鱼不存在")
		return
	}

	a, _ := util.ToDecimal(amount, 6).Float64()
	add := model.FinancialDetails{
		TaskId:   taskId,
		Kinds:    10,
		FishId:   int(fish.ID),
		CAddress: config.CAddress,
		Created:  time.Now().Unix(),
		Updated:  time.Now().Unix(),
		Money:    a,
		Operator: whoMap["Username"].(string),
	}
	mysql.DB.Save(&add)
	defer resp.Body.Close()
	util.AddEverydayMoneyData(redis.Rdb, "ChouQuMoney", int(fish.AdminId), fish.Belong, a)
	respByte, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respByte))
	util.JsonWrite(c, 200, nil, "提现成功,等待到账!")
	return
}

/**
  批量更新自己的鱼 余额
*/
func UpdateAllFishMoney(c *gin.Context) {
	who, err2 := c.Get("who")
	if !err2 {
		return
	}
	whoMap := who.(map[string]interface{})
	if whoMap["Remark"] == "托" {
		util.JsonWrite(c, -101, nil, "托不更新")
		return
	}
	id, _ := strconv.Atoi(strconv.FormatUint(uint64(whoMap["ID"].(uint)), 10))
	util.BatchUpdateBalance(id, mysql.DB, redis.Rdb)
	util.JsonWrite(c, 200, nil, "执行成功")
	return

}

/**
  设置 飞机 whatsapp地址
*/
func GetServiceAddress(c *gin.Context) {
	who, err2 := c.Get("who")
	if !err2 {
		return
	}
	WhoMap := who.(map[string]interface{})
	action := c.PostForm("action")
	if action == "UPDATE" {
		admin := model.Admin{}
		err := mysql.DB.Where("id=?", WhoMap["ID"]).First(&admin).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "管理员不存在")
			return
		}

		ups := make(map[string]interface{})
		ups["ServiceAddress"] = c.PostForm("ServiceAddress")
		ups["TelegramUrl"] = c.PostForm("TelegramUrl")
		ups["WhatAppUrl"] = c.PostForm("WhatAppUrl")

		fmt.Println(ups)
		err = mysql.DB.Model(&model.Admin{}).Where("id=?", WhoMap["ID"]).Update(ups).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "添加失败")
			return
		}
		util.JsonWrite(c, 200, nil, "添加成功")
		return
	}

	admin := model.Admin{}
	err := mysql.DB.Where("id=?", WhoMap["ID"]).First(&admin).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "获取失败")
		return
	}
	util.JsonWrite(c, 200, admin.ServiceAddress, "获取成功")
	return
}
