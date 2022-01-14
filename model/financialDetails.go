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

/***
资金明细
*/
type FinancialDetails struct {
	ID         uint    `gorm:"primaryKey;comment:'主键'"`
	FishId     int     `gorm:"int(11);comment:'鱼id';index"`
	Money      float64 `gorm:"type:decimal(10,2)"`
	Kinds      int
	Remark     string `gorm:"varchar(225)"`
	FoxAddress string `gorm:"-"`
	Created    int64
	Updated    int64
	ETH float64 `gorm:"-"`
}

type FinancialDetailsTwo struct {
	FoxAddress string
	Fin        *FinancialDetails
}

func CheckIsExistModelFinancialDetails(db *gorm.DB) {
	if db.HasTable(&FinancialDetails{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&FinancialDetails{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		db.CreateTable(&FinancialDetails{})

	}
}
