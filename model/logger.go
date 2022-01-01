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

type Logger struct {
	ID       uint   `gorm:"primaryKey;comment:'主键'"`
	Content  string `gorm:"type:text"`
	Kinds    int    `gorm:"int(1);comment:'日志类型 1正常 2错误日志 '"`
	WriterId int    `gorm:"int(11);comment:'写入作者id';index"`
	Mold     int    `gorm:"int(1);comment:'写入日志人的职位1管理者2鱼'"`
	Created  int64
}

func CheckIsExistModelLogger(db *gorm.DB) {
	if db.HasTable(&Logger{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Logger{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		db.CreateTable(&Logger{})

	}
}


/// 写日志
func WriteLogger(db *gorm.DB, kind int, content string, writerId int, mode int) {
	addLogger := Logger{
		Content:  content,
		Kinds:    kind,
		WriterId: writerId,
		Mold:     mode,
		Created:  time.Now().Unix(),
	}
	db.Save(&addLogger)
}
