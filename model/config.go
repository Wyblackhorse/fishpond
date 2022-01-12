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
	ID        uint   `gorm:"primaryKey;comment:'主键'"`
	BAddress  string `gorm:"varchar(225)"`
	BKey      string `gorm:"varchar(225)"`
	BMnemonic string `gorm:"varchar(225)"`
	CAddress  string `gorm:"varchar(225)"`
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
