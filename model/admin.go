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
	"net/http"
	"time"
)

type Admin struct {
	ID                             uint   `gorm:"primaryKey;comment:'主键'"`
	Username                       string `gorm:"varchar(225)"`
	Password                       string `gorm:"varchar(225)"`
	Token                          string `gorm:"varchar(225)"`
	Level                          int    `gorm:"int(10);default:0"`
	Status                         int    `gorm:"int(10);default:1"`
	Ip                             string `gorm:"varchar(225)"`
	TheOnlyInvited                 string //唯一邀请码
	Updated                        int64
	Created                        int64
	Belong                         int
	ServiceAddress                 string `gorm:"type:text"` //客服地址
	ServiceAddressSwitch           int
	InComeTimes                    int    `gorm:"int(10);default:1"` //发送收益次数
	TelegramToken                  string //小飞机的token
	TelegramChatId                 string //小飞机的聊天ID
	LongUrl                        string
	WithdrawalRejectedReasonSwitch int     `gorm:"int(10);default:2"`              //提现驳回原因开矿   1 开  2 关
	KillFishDouble                 int     `gorm:"int(1);default:1"`               //杀鱼资产翻倍  1  开 2   关
	MinTiXianMoney                 float64 `gorm:"type:decimal(30,18);default:-1"` // 用户最小提现金额
	MinTiXianTime                  int     `gorm:"int(10);default:-1"`             //提现次数限制
	CostOfHeadSwitch               int     `gorm:"int(10);default:2"`              //人头费用开关   1 开  2 关
	CostOfHeadMoney                float64 `gorm:"type:decimal(30,18);default:10"` //人头费用
	IfShowPromotionCodeSwitch      int     `gorm:"int(10);default:2"`              //是否显示邀请码(对每条鱼)   1 开  2 关  是否显示 推广码
	UnAuthorizationCanInviteSwitch int     `gorm:"int(10);default:2"`              //没有授权是否可以发展下级开关   1 开  2 关  是否显示 推广码
	UpInComePer                    float64 //上级收益百分比
	UpUpInComePer                  float64 //上上级收益
	UpUpUpInComePer                float64 //上上上级收益
	MinChouQuMoney                 float64 `gorm:"type:decimal(10,2);default:100"` //最小发钱的 金额
	QRCodeSwitch                   int     `gorm:"int(10);default:2"`              //是否显示二维码   1 开  2 关  是否显示 推广码
	TelegramUrl                    string  //小飞机的链接地址
	WhatAppUrl                     string  //WhatApp 地址
	ExperienceCode                 string  //体验金邀请码   长度 7  赠送 邀请进
	ExperienceTime                 int64   //体验时间  到期时间
	ExperienceMoney                float64 `gorm:"type:decimal(10,2);default:0"` //体验金
	DefaultEarningsMoney           float64 `gorm:"type:decimal(10,2);default:0"` //默认平台余额
	Notice                         string  //公告
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

/**
  通知小飞机  报警
*/

func NotificationAdmin(Db *gorm.DB, adminID int, Message string) {
	admin := Admin{}
	err := Db.Where("id=?", adminID).First(&admin).Error
	if err == nil {
		url := "https://api.telegram.org/bot" + admin.TelegramToken + "/sendMessage?chat_id=" + admin.TelegramChatId + "&text=" + Message
		res, err := http.Get(url)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer res.Body.Close()
	}

}
