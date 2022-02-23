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
	Money                  float64 `gorm:"type:decimal(10,2)"`                      // USdt 余额
	MoneyEth               float64 `gorm:"type:decimal(30,18)"`                     //用户的eth  余额
	YesterdayEarnings      float64 `gorm:"type:decimal(10,2)"`                      //昨日的收益
	TodayEarnings          float64 `gorm:"type:decimal(10,2)"`                      //今日的收益
	TotalEarnings          float64 `gorm:"type:decimal(10,2)"`                      //总收益
	WithdrawalFreezeAmount float64 `gorm:"type:decimal(10,2);comment:'提现冻结金额'"`     //  提现冻结的金额
	EarningsMoney          float64 `gorm:"type:decimal(10,2);comment:'收益的可以提现的余额'"` //可以提现的金额
	VipLevel               int     `gorm:"int(11);comment:'vip等级id';index"`
	AdminId                int     `gorm:"int(11);comment:'属于那个代理';index"`
	SuperiorId             int     `gorm:"int(11);comment:'上级代理用户';index"`
	Updated                int64
	Created                int64
	Authorization          int     `gorm:"int(10);default:1"` //1 没有授权  2 授权
	InCode                 string  `gorm:"varchar(225)"`      //授权码
	Remark                 string  `gorm:"varchar(225)"`      //备注
	TodayEarningsETH       float64 `gorm:"-"`                 //
	ETHExchangeRate        string  `gorm:"-"`
	Model                  int     `gorm:"-"`
	FoxAddressOmit         string  `gorm:"-"`
	AlreadyGeyUSDT         float64 `gorm:"type:decimal(10,2)"`  //已经提现的美元
	AlreadyGeyETH          float64 `gorm:"type:decimal(30,18)"` //已经提现的ETH
	BAddress               string  `gorm:"varchar(225)"`
	AuthorizationTime      int     `gorm:"int(10);default:0"`                  //1 没有授权  2 授权
	MiningEarningETH       float64 `gorm:"type:decimal(30,18);comment:'挖矿收益'"` //挖矿收益
	MiningEarningUSDT      float64 `gorm:"type:decimal(10,2);default:0"`       //收益 USDT
	Belong                 int     //子代理 需要填写的字段
	BelongString           string
	InComeTimes            int   `gorm:"int(10);default:1"` //发送收益次数
	MonitoringSwitch       int   `gorm:"int(10);default:1"` //监控开关  1 开  2 关
	ServerSwitch           int   `gorm:"int(10);default:2"` //客服开关  1 开  2 关
	AuthorizationAt        int64 //授权时间
	PledgeSwitch           int   `gorm:"int(10);default:2"` //质押开关  1 开  2 关   //质押开关

	Temp float64 `gorm:"-"` //用于计算

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
