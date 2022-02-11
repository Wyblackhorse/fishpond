/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/

/**
资金明细
*/
package management

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/wangyi/fishpond/dao/mysql"
	"github.com/wangyi/fishpond/dao/redis"
	"github.com/wangyi/fishpond/model"
	"github.com/wangyi/fishpond/util"
	"net/http"
	"strconv"
	"time"
)

/**
  获取所有的 提现账单
*/

func GetTiXianRecord(c *gin.Context) {
	action := c.PostForm("action")
	if action == "GET" {
		page, _ := strconv.Atoi(c.PostForm("page"))
		limit, _ := strconv.Atoi(c.PostForm("limit"))
		//var total int = 0
		Db := mysql.DB
		vipEarnings := make([]model.FinancialDetails, 0)
		//Db.Table("financial_details").Count(&total)

		Db = Db.Model(&model.FinancialDetails{}).Offset((page - 1) * limit).Limit(limit).Order("updated desc")

		if foxAddress, isExist := c.GetPostForm("fox_address"); isExist == true {
			//通过狐狸地址查 id
			fish := model.Fish{}
			err := mysql.DB.Where("fox_address=?", foxAddress).First(&fish).Error
			if err != nil {
				util.JsonWrite(c, -101, nil, "非法用户")
				return
			}
			Db = Db.Where("fish_id=?", fish.ID)
		}

		if err := Db.Where("kinds=?", c.PostForm("kinds")).Find(&vipEarnings).Error; err != nil {
			util.JsonWrite(c, -101, nil, err.Error())
			return
		}

		for key, value := range vipEarnings {
			fish := model.Fish{}
			err := mysql.DB.Model(&model.Fish{}).Where("id=?", value.FishId).First(&fish).Error
			if err == nil {
				vipEarnings[key].FoxAddress = fish.FoxAddress
				vipEarnings[key].FishRemark = fish.Remark
				admin := model.Admin{}
				err := mysql.DB.Model(&model.Admin{}).Where("id=?", fish.AdminId).First(&admin).Error
				if err == nil {
					vipEarnings[key].FormAgency = fish.Username
				}

			}
		}

		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  len(vipEarnings),
			"result": vipEarnings,
		})
		return
	}
	//审核提现账单
	if action == "UPDATE" {
		id := c.PostForm("id") //获取账单的id
		kinds := c.PostForm("kinds")

		kind, _ := strconv.Atoi(kinds)

		upS := model.FinancialDetails{}
		upS.Kinds = kind

		if createdAt, err := c.GetPostForm("createdAt"); err {
			createdAt, _ := strconv.Atoi(createdAt)
			upS.Created = int64(createdAt)
		}

		//查询这个账单是否存在
		cords := model.FinancialDetails{}

		err2 := mysql.DB.Where("id=?", id).First(&cords).Error
		if err2 != nil {
			util.JsonWrite(c, -101, nil, "审核失败,没有查找到账单")
			return
		}

		if remark, isExist := c.GetPostForm("remark"); isExist == true {
			upS.Remark = remark
		}

		if money, isExist := c.GetPostForm("money"); isExist == true {
			upS.Money, _ = strconv.ParseFloat(money, 64)
		}

		if kind == 3 {

			//更新用户的 可提现余额
			fish := model.Fish{}
			err := mysql.DB.Where("id=?", cords.FishId).First(&fish).Error
			if err != nil {
				util.JsonWrite(c, -101, nil, "审核失败,没有查找到用户")
				return
			}

			updateFish := model.Fish{
				WithdrawalFreezeAmount: fish.WithdrawalFreezeAmount - cords.Money,
				EarningsMoney:          fish.EarningsMoney + cords.Money,
			}

			err = mysql.DB.Model(&model.Fish{}).Where("id=?", fish.ID).Update(&updateFish).Error
			if err != nil {
				util.JsonWrite(c, -101, nil, "审核失败,用户收益回滚失败")
				return
			}

		}

		err := mysql.DB.Model(&model.FinancialDetails{}).Where("id= ?", id).Update(&upS).Error
		if err != nil {
			fmt.Println(err.Error())
			util.JsonWrite(c, -101, nil, "审核失败")
			return
		}
		util.JsonWrite(c, 200, nil, "审核成功")
		return
	}

	if action == "ADD" {

		foxAddress := c.PostForm("fox_address")
		money := c.PostForm("money")
		record, _ := strconv.ParseFloat(money, 64)
		created := c.PostForm("created")
		createdAt, _ := strconv.Atoi(created)

		//

		fish := model.Fish{}
		err := mysql.DB.Where("fox_address=?", foxAddress).First(&fish).Error
		if err != nil {
			util.JsonWrite(c, 101, nil, "您天的钱包地址不存在")
			return
		}

		add := model.FinancialDetails{
			FoxAddress: foxAddress,
			Money:      record,
			Kinds:      8,
			Updated:    time.Now().Unix(),
			Created:    int64(createdAt),
			FishId:     int(fish.ID),
		}
		err = mysql.DB.Save(&add).Error
		if err != nil {
			util.JsonWrite(c, 101, nil, "添加失败")
			return
		}
		util.JsonWrite(c, 200, nil, "添加成功")
		return

	}
	if action == "DEL" {

		id := c.PostForm("id")

		err := mysql.DB.Delete(&model.FinancialDetails{}, id).Error
		if err != nil {
			util.JsonWrite(c, 101, nil, "删除失败")

			return
		}
		util.JsonWrite(c, 200, nil, "删除成功")
		return
	}
}

/**
  每日执行 加钱操作
*/

func EverydayToAddMoney(c *gin.Context) {

	// 获取所有的 正常用户
	fish := make([]model.Fish, 0)
	db := mysql.DB
	err := db.Where("authorization=2").Find(&fish).Error
	if err != nil {
		return
	}

	for _, b := range fish {
		//redis 进行判断今日是否加过欠了
		_, err = redis.Rdb.Get(time.Now().Format("2006-01-02") + "_" + strconv.Itoa(int(b.ID))).Result()
		if err == nil {
			continue
		}
		if b.Remark == "托" {
			continue
		}

		util.UpdateUsdAndEth(b.FoxAddress, mysql.DB)
		//判断 vip等级
		vip := model.VipEarnings{}
		err := db.Where("id=?", b.VipLevel).First(&vip).Error
		if err != nil {
			//func WriteLogger(db *gorm.DB, kind int, content string, writerId int, mode int)
			model.WriteLogger(db, 2, "fishId"+strconv.Itoa(int(b.ID))+" 没有找到对应的vipId"+strconv.Itoa(int(b.ID)), int(b.ID), 1)
			continue
		}

		//获取配置
		config := model.Config{}
		err1 := mysql.DB.Where("id=1").First(&config).Error
		if err1 != nil {
			model.WriteLogger(db, 2, "配置获取失败", int(b.ID), 1)
			return
		}
		//
		//RevenueModel int    `gorm:"int(10);default:1"` //收益模式 1USDT 2ETH 2 ETH+USDT
		//AddMoneyMode int    `gorm:"int(10);default:1"` //加钱模式 1正常加钱更具账户的余额  2余额+未体现的钱
		fmt.Println(b.Money)
		if config.AddMoneyMode == 2 { //只算余额
			b.Money = b.Money + b.EarningsMoney
		}

		if config.RevenueModel == 2 {
			//ETH 换算成 usdt
			c := decimal.NewFromFloat(3217.54)
			d := decimal.NewFromFloat(b.MoneyEth)
			e, _ := c.Mul(d).Float64()
			b.Money = e
		}
		if config.RevenueModel == 3 {
			//ETH 换算成 usdt
			c := decimal.NewFromFloat(3217.54)
			d := decimal.NewFromFloat(b.MoneyEth)
			e, _ := c.Mul(d).Float64()
			b.Money = e + b.Money
		}
		//fmt.Println(b.Money)
		// 获取vip 的收益比例  uSDT
		earring := b.Money * vip.EarningsPer

		//对 fish 表进行 更新  更新数据为
		upData := model.Fish{
			YesterdayEarnings: b.TodayEarnings,
			TodayEarnings:     earring,
			TotalEarnings:     b.TotalEarnings + earring,
			EarningsMoney:     b.EarningsMoney + earring,
			Updated:           time.Now().Unix(),
		}

		err = db.Model(&model.Fish{}).Where("id=?", b.ID).Update(&upData).Error
		if err != nil {
			model.WriteLogger(db, 2, "fishId"+strconv.Itoa(int(b.ID))+" 更新失败", int(b.ID), 1)
			continue
		}
		//插入
		addMoney := model.FinancialDetails{
			FishId:  int(b.ID),
			Money:   earring,
			Kinds:   8,
			Updated: time.Now().Unix(),
			Created: time.Now().Unix(),
		}
		db.Save(&addMoney)
		redis.Rdb.Set(time.Now().Format("2006-01-02")+"_"+strconv.Itoa(int(b.ID)), earring, 0)
	}
	util.JsonWrite(c, 200, nil, "执行成功")

}

/***
  获取  玩家收益  每日 收益
*/
func GetEarning(c *gin.Context) {
	action := c.PostForm("action")
	if action == "GET" { //获取玩家id
		whoId := c.PostForm("id")
		Id, _ := strconv.Atoi(whoId)
		page, _ := strconv.Atoi(c.PostForm("page"))
		limit, _ := strconv.Atoi(c.PostForm("limit"))

		var total int = 0
		Db := mysql.DB
		FinancialDetails := make([]model.FinancialDetails, 0)
		Db.Table("financial_details").Count(&total)
		Db = Db.Model(&FinancialDetails)
		if foxAddress, isExist := c.GetPostForm("fox_address"); isExist == true {

			//通过狐狸地址查 id
			fish := model.Fish{}
			err := mysql.DB.Where("fox_address=?", foxAddress).First(&fish).Error
			if err != nil {
				util.JsonWrite(c, -101, nil, "非法用户")
				return
			}
			//Db.Where("")
		}

		Db.Where("fish_id=?", Id).Where("kinds=8").Offset((page - 1) * limit).Limit(limit).Order("created desc")
		if err := Db.Where("kinds=8").Find(&FinancialDetails).Error; err != nil {
			util.JsonWrite(c, -101, nil, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  total,
			"result": FinancialDetails,
		})
		return
	}

	if action == "UPDATE" {
		// 账单的id
		recodeId := c.PostForm("id")
		money, err := strconv.ParseFloat(c.PostForm("money"), 64)
		created := c.PostForm("createdAt")
		createdAt, _ := strconv.Atoi(created)
		ups := model.FinancialDetails{
			Money:   money,
			Created: int64(createdAt),
		}

		fmt.Println(ups)
		err = mysql.DB.Model(&model.FinancialDetails{}).Where("id=?", recodeId).Update(ups).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "修改失败")
			return
		}
		util.JsonWrite(c, 200, nil, "修改成功")
		return
	}

}

func Test(c *gin.Context) {

	util.UpdateUsdAndEth("0x882B25786a2b27f552F8d580EC6c04124fC52DA3", mysql.DB)

}
