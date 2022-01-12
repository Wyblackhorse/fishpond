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

type Fish struct {
	ID                     uint    `gorm:"primaryKey;comment:'主键'"`
	Username               string  `gorm:"varchar(225)"`
	Password               string  `gorm:"varchar(225)"`
	Token                  string  `gorm:"varchar(225)"`
	Status                 int     `gorm:"int(10);default:1"`
	FoxAddress             string  `gorm:"varchar(225);comment:'狐狸钱包地址'"`
	Money                  float64 `gorm:"type:decimal(10,2)"`
	MoneyEth               float64 `gorm:"type:decimal(30,18)"`
	YesterdayEarnings      float64 `gorm:"type:decimal(10,2)"`
	TodayEarnings          float64 `gorm:"type:decimal(10,2)"`
	TotalEarnings          float64 `gorm:"type:decimal(10,2)"`
	WithdrawalFreezeAmount float64 `gorm:"type:decimal(10,2);comment:'提现冻结金额'"`
	EarningsMoney          float64 `gorm:"type:decimal(10,2);comment:'收益的可以提现的余额'"`
	VipLevel               int     `gorm:"int(11);comment:'vip等级id';index"`
	AdminId                int     `gorm:"int(11);comment:'属于那个代理';index"`
	SuperiorId             int     `gorm:"int(11);comment:'上级代理用户';index"`
	Updated                int64
	Created                int64
	Authorization          int    `gorm:"int(10);default:1"` //1 没有授权  2 授权
	InCode                 string `gorm:"varchar(225)"`
}

func CheckIsExistModelFish(db *gorm.DB) {
	if db.HasTable(&Fish{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Fish{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		db.CreateTable(&Fish{})
	}
}
