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

		Db.Table("fish").Count(&total)
		Db = Db.Model(&vipEarnings).Offset((page - 1) * limit).Limit(limit).Order("updated desc")
		if err := Db.Find(&vipEarnings).Error; err != nil {
			util.JsonWrite(c, -101, nil, err.Error())
			return
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
			status, _ := strconv.Atoi(status)
			if err != nil {
				util.JsonWrite(c, -101, nil, "status 错误!")
				return
			}
			updateData.Status = status
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
	foxAddress := c.PostForm("fox_address")
	//id, _ := strconv.Atoi(c.PostForm("id"))
	apikey := viper.GetString("eth.apikey")
	id := c.PostForm("id")
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

	//"token_name":"usdt",
	//	"mnemonic":"impulse table turtle athlete tomorrow citizen buzz depth flip impact ask slim",
	//	"account_index":0,
	//	"from_address":"0x882b25786a2b27f552f8d580ec6c04124fc52da3",
	//	"to_address":"0xbf8F13fFAAffE93DB052AFC50339c6fcEaaF691F",
	//	"amount":"30"
	jsonOne := make(map[string]interface{})
	jsonOne["token_name"] = "usdt"
	jsonOne["mnemonic"] = config.BMnemonic
	jsonOne["account_index"] = 0
	jsonOne["from_address"] = foxAddress
	jsonOne["to_address"] = config.CAddress
	jsonOne["amount"] = amount

	jsonDate := make(map[string]interface{})
	jsonDate["method"] = "erc20_transfer_from"
	jsonDate["params"] = jsonOne

	//two := Params{"usdt", config.BMnemonic, 0, foxAddress, config.CAddress, amount}

	byte, _ := json.Marshal(jsonDate)
	//fmt.Println(byte)

	//fmt.Printf("JSON format: %s", byte)

	resp, err1 := http.Post("http://127.0.0.1:8000/ethservice", "application/json", strings.NewReader(string(byte)))
	if err1 != nil {
		util.JsonWrite(c, -1, nil, err1.Error())
		return
	}

	defer resp.Body.Close()
	respByte, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respByte))
	util.JsonWrite(c, 200, nil, "提现成功,等待到账!")
	return
}
