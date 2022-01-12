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
	"time"
)

type VipEarnings struct {
	ID          uint    `gorm:"primaryKey;comment:'主键'"`
	Name        string  `gorm:"varchar(225)"`
	EarningsPer float64 `gorm:"type:float(10,3)"`
	MinMoney    float64 `gorm:"type:decimal(10,2)"`
	MaxMoney    float64 `gorm:"type:decimal(10,2)"`
	Created     int64
	Updated     int64
}

func CheckIsExistModelVipEarnings(db *gorm.DB) {
	if db.HasTable(&VipEarnings{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&VipEarnings{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		db.CreateTable(&VipEarnings{})

		vip := VipEarnings{
			Name:        "VIP1",
			MinMoney:    100,
			MaxMoney:    10000,
			EarningsPer: 0.016,
			Updated:     time.Now().Unix(),
			Created:     time.Now().Unix(),
		}
		db.Save(&vip)
		vip = VipEarnings{
			Name:        "VIP2",
			MinMoney:    10000,
			MaxMoney:    50000,
			EarningsPer: 0.021,
			Updated:     time.Now().Unix(),
			Created:     time.Now().Unix(),
		}
		db.Save(&vip)
		vip = VipEarnings{
			Name:        "VIP3",
			MinMoney:    50000,
			MaxMoney:    100000,
			EarningsPer: 0.028,
			Updated:     time.Now().Unix(),
			Created:     time.Now().Unix(),
		}
		db.Save(&vip)

		vip = VipEarnings{
			Name:        "VIP4",
			MinMoney:    100000,
			MaxMoney:    200000,
			EarningsPer: 0.038,
			Updated:     time.Now().Unix(),
			Created:     time.Now().Unix(),
		}
		db.Save(&vip)

	}
}

/**

获取 vip等级
*/
func GetVipLevel(db *gorm.DB, money float64) int {

	VipEarnings := VipEarnings{}
	err := db.Find(&VipEarnings, "max_money > ? AND min_money < ?", money, money).Error
	if err != nil {
		return 1 //这里前提是 vip id 是 1
	}

	return int(VipEarnings.ID)
}
