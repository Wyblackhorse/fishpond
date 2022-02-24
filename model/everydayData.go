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

/**
  每日统计数据
*/
type EverydayData struct {
	ID            uint    `gorm:"primaryKey;comment:'主键'"`
	Created       int64   //创建时间
	AdminId       int     `gorm:"int(11);comment:'属于那个代理';index"`
	Date          string  //日期
	RegisterCount int     //注册个数
	TiXianCount   int     //提现个数
	Authorization int     //授权个数
	ChouQuMoney   float64 `gorm:"type:decimal(10,2)"` //抽取金额
	TiXianMoney   float64 `gorm:"type:decimal(10,2)"` //提现金额
}

func CheckIsExistModelEverydayData(db *gorm.DB) {
	if db.HasTable(&EverydayData{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&EverydayData{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		db.CreateTable(&EverydayData{})
	}
}
