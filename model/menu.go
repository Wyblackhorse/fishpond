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

type Menu struct {
	ID      uint   `gorm:"primaryKey;comment:'主键'"`
	Belong  string `gorm:"varchar(225)"`
	Name    string `gorm:"varchar(225)"`
	Status  int    `gorm:"int(10);default:1"`
	Level   int
	Created int64
}

func CheckIsExistModelMenu(db *gorm.DB) {
	if db.HasTable(&Menu{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Menu{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&Menu{}).Error
		if err == nil {
			//创建成功  这里就插入超级管理员
			addMenu := Menu{
				Name:    "控制台",
				Level:   0,
				Status:  1,
				Created: time.Now().Unix(),
				Belong:  "0",
			}
			err := db.Save(&addMenu).Error
			if err != nil {
				fmt.Println("插入失败")
			}

			addMenu1 := Menu{
				Name:    "鱼塘管理",
				Status:  1,
				Created: time.Now().Unix(),
				Belong:  "0",
				Level:   0,
			}
			db.Save(&addMenu1)

			//db.Save(&Menu{Name: "活跃的鱼",Level: })
			
			addMenu2 := Menu{
				Name:    "Vip权限管理",
				Status:  1,
				Created: time.Now().Unix(),
				Belong:  "0",
			}
			db.Save(&addMenu2)

			addMenu3 := Menu{
				Name:    "代理管理",
				Status:  1,
				Created: time.Now().Unix(),
				Belong:  "0",
			}
			db.Save(&addMenu3)

			addMenu4 := Menu{
				Name:    "系统设置",
				Status:  1,
				Created: time.Now().Unix(),
				Belong:  "0",
			}
			db.Save(&addMenu4)

		}

	}
}

func GetMenus(db *gorm.DB, level string) []Menu {
	menu := make([]Menu, 0)

	db.Where("belong LIKE ?", "%"+level+"%").Find(&menu)
	return menu

}
