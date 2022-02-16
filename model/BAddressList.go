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

type BAddressList struct {
	ID       uint   `gorm:"primaryKey;comment:'主键'"`
	BAddress string `gorm:"varchar(225)"`
	BKey     string `gorm:"varchar(225)"`
}

func CheckIsExistModeBAddressList(db *gorm.DB) {
	if db.HasTable(&BAddressList{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&BAddressList{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		db.CreateTable(&BAddressList{})

	}
}

/**
  添加新的B地址
*/
func AddBAddressList(DB *gorm.DB, B string, key string) {
	//
	list := BAddressList{}
	err := DB.Where("b_address=?", B).First(&list).Error
	if err != nil {
		//不存在
		add := BAddressList{
			BAddress: B,
			BKey:     key,
		}
		DB.Save(&add)
		return
	}
	//存在
	DB.Model(&BAddressList{}).Where("id=?", list.ID).Update(&BAddressList{
		BKey: key,
	})

}
