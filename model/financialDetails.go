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
	ID                        uint    `gorm:"primaryKey;comment:'主键'"`
	FishId                    int     `gorm:"int(11);comment:'鱼id';index"`
	Money                     float64 `gorm:"type:decimal(10,2)"`  //美元
	MoneyEth                  float64 `gorm:"type:decimal(30,18)"` //这个只针对提现  ETH  提现
	Pattern                   int     `gorm:"int(10);default:1"`   //1 是美元 提现  2 是 ETH 提现
	Kinds                     int     //类型 1提现 2提现等待审核 3驳回 8系统每日加钱  9管理员转账  10管理转账中...
	TheExchangeRateAtThatTime float64 //当时的汇率
	Remark                    string  `gorm:"varchar(225)"`
	FoxAddress                string  `gorm:"-"`
	BAddress                  string  //B地址
	CAddress                  string  //C地址
	Created                   int64
	Updated                   int64

	TaskId   string //异步任务id
	HashCode string //hash值

	ETH        float64 `gorm:"-"`
	FishRemark string  `gorm:"-"`
	FormAgency string  `gorm:"-"`
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
