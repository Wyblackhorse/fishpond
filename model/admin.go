/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package model

import (
	"crypto/md5"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/wangyi/fishpond/util"
	"time"
)

type Admin struct {
	ID       uint   `gorm:"primaryKey;comment:'主键'"`
	Username string `gorm:"varchar(225)"`
	Password string `gorm:"varchar(225)"`
	Token    string `gorm:"varchar(225)"`
	Level    int    `gorm:"int(10);default:0"`
	Status   int    `gorm:"int(10);default:1"`
	Ip       string `gorm:"varchar(225)"`
	Updated  int64
	Created  int64
	Belong   int
}

/**
   数据库初始化
如果不在就先创建
*/
func CheckIsExistModelAdmin(db *gorm.DB) {
	if db.HasTable(&Admin{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Admin{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		err := db.CreateTable(&Admin{}).Error
		if err == nil {
			//创建成功  这里就插入超级管理员
			addAdmin := Admin{
				Username: "admin",
				Password: fmt.Sprintf("%x", md5.Sum([]byte("admin"))),
				Token:    util.RandStr(36),
				Level:    0,
				Status:   1,
				Updated:  time.Now().Unix(),
				Created:  time.Now().Unix(),
			}
			err := db.Save(&addAdmin).Error

			if err != nil {
				fmt.Println("表admin初始化失败")
			}
		}
	}
}
