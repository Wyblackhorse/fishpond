/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

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
	LowCanKillFishMoney     float64 `gorm:"int(10);default:50"`                        //美元

}

func CheckIsExistModelConfig(db *gorm.DB) {
	if db.HasTable(&Config{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Config{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&Config{}).Error
		if err == nil {
			//创建成功  这里就插入超级管理员
			Config := Config{
				ID: 1,
			}
			err := db.Save(&Config).Error

			if err != nil {
				fmt.Println("表admin初始化失败")
			}
		}
	}
}
