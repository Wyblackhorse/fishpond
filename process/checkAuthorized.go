/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package process

import (
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
		err := Db.Where("authorization=1").Find(&fish).Error
		if err == nil {
			for _, value := range fish {
				config := model.Config{}
				err := Db.Where("id=1").First(&config).Error
				if err == nil {
					util.ChekAuthorizedFoxAddress(value.FoxAddress, apiKey, config.BAddress, Db)
				}
			}
		}

		time.Sleep(3600 * time.Second)

	}

}
