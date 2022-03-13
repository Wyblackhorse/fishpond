/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package process

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
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
				time.Sleep(5 * time.Second)
			}
		}
		time.Sleep(600 * time.Second)
	}
}
