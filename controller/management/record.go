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
		var total int

		//判断是否有总代的 id
		if _, isExist := c.GetPostForm("adminId"); isExist == true {
			if _, isExist := c.GetPostForm("SonAdminId"); isExist == true {
				//子代也存在
				adminId, _ := strconv.Atoi(c.PostForm("SonAdminId"))
				Db = Db.Table("financial_details").Joins("left join fish on fish.id=financial_details.fish_id ").Where("fish.admin_id= ?", adminId)
			} else {
				adminId, _ := strconv.Atoi(c.PostForm("adminId"))
				Db = Db.Table("financial_details").Joins("left join fish on fish.id=financial_details.fish_id ").Where("fish.belong= ?", adminId)
			}
		} else {
			//查询所有的鱼
			Db = Db.Table("financial_details").Joins("left join fish on fish.id=financial_details.fish_id ")
		}

		//sonAdmins := make([]model.Admin, 0)
		//err := mysql.DB.Where("belong=?", whoMap["ID"]).Find(&sonAdmins).Error
		//if err != nil {
		//	util.JsonWrite(c, -101, nil, "查询失败")
		//	return
		//}

		//if _, isExist := c.GetPostForm("adminId"); isExist == true {
		//	adminId, _ := strconv.Atoi(c.PostForm("adminId"))
		//	Db = Db.Table("financial_details").Joins("left join fish on fish.id=financial_details.fish_id ").Where("fish.admin_id= ?", adminId)
		//} else {
		//	var BString []string
		//	BString = append(BString, strconv.FormatUint(uint64(whoMap["ID"].(uint)), 10))
		//	for _, v := range sonAdmins {
		//		BString = append(BString, strconv.Itoa(int(v.ID)))
		//	}
		//	Db = Db.Table("financial_details").Joins("left join fish on fish.id=financial_details.fish_id ").Where("fish.admin_id IN (?)", BString)
		//}

		if _, isExist := c.GetPostForm("tuo"); isExist == true {
			Db = Db.Where("fish.remark!=?", "托")
		}

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
		if pattern, isExist := c.GetPostForm("pattern"); isExist == true {
			Db = Db.Where("pattern=?", pattern)
		}

		if start, isExist := c.GetPostForm("start"); isExist == true {
			if end, isExist := c.GetPostForm("end"); isExist == true {
				Db = Db.Where("financial_details.created < ? AND financial_details.created > ? ", end, start)
			}
		}

		Db = Db.Where("kinds=?", c.PostForm("kinds"))
		Db.Where("kinds=?", c.PostForm("kinds")).Offset((page - 1) * limit).Limit(limit).Order("updated desc").Find(&vipEarnings)
		if err := Db.Offset((page - 1) * limit).Limit(limit).Order("updated desc").Find(&vipEarnings).Error; err != nil {
			util.JsonWrite(c, -101, nil, err.Error())
			return
		}

		Db.Count(&total)
		for key, value := range vipEarnings {
			fish := model.Fish{}
			err := mysql.DB.Model(&model.Fish{}).Where("id=?", value.FishId).First(&fish).Error
			if err == nil {
				vipEarnings[key].FoxAddress = fish.FoxAddress
				vipEarnings[key].FishRemark = fish.Remark
				admin := model.Admin{}
				err := mysql.DB.Model(&model.Admin{}).Where("id=?", fish.AdminId).First(&admin).Error
				if err == nil {
					vipEarnings[key].FormAgency = admin.Username
				}

			}
		}

		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  total,
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
  每日执行 加钱操作  这个是 一个总的
*/

func EverydayToAddMoney(c *gin.Context) {
	// 获取所有的 正常用户
	fish := make([]model.Fish, 0)
	db := mysql.DB
	err := db.Where("authorization=2 or remark=? or no_proceeds_are_authorized_switch=1", "托").Find(&fish).Error
	if err != nil {
		fmt.Printf(err.Error())
		//`gorm:"int(1);comment:'日志类型 1正常 2错误日志 '"`
		model.WriteLogger(db, 2, "EverydayToAddMoney 错误:"+err.Error(), 0, 1)
		return
	}

	for _, b := range fish {
		model.WriteLogger(db, 2, "准备发放收益", int(b.ID), 2)
		TimesOne, err := redis.Rdb.Get(time.Now().Format("2006-01-02") + "_" + strconv.Itoa(int(b.ID))).Result()
		if err == nil && b.InComeTimes == 1 { //有数据 但是针对 每天只发 一次收益的 玩家停止
			continue
		}
		if b.InComeTimes == 2 { //玩家发每天发两次收益
			ppp, _ := strconv.Atoi(TimesOne)
			if ppp == 2 {
				model.WriteLogger(db, 2, "EverydayToAddMoney  今日已经运行结束", int(b.ID), 2)
				continue
			}
		}
		if b.Remark != "托" {
			util.UpdateUsdAndEth(b.FoxAddress, mysql.DB, b.Money, int(b.ID), b.AdminId, b.Remark, redis.Rdb)
		}
		//获取配置
		config := model.Config{}
		err1 := mysql.DB.Where("id=1").First(&config).Error
		if err1 != nil {
			model.WriteLogger(db, 2, "配置获取失败", int(b.ID), 1)
			return
		}
		ethHl, _ := redis.Rdb.Get("ETHTOUSDT").Result()
		ETH2, _ := strconv.ParseFloat(ethHl, 64) ////收益模式 1USDT 2ETH 2 ETH+USDT
		if config.RevenueModel == 2 {
			//ETH 换算成 usd
			c := decimal.NewFromFloat(ETH2)
			d := decimal.NewFromFloat(b.MoneyEth)
			e, _ := c.Mul(d).Float64()
			b.Money = e
		}
		//判断  奖励金是否到期
		if b.ExperienceMoney > 0 && b.ExpirationTime > time.Now().Unix() { //奖励金 必须大于0 并且没有  过期
			b.Money = b.Money + b.ExperienceMoney
		}

		if b.Balance > 0 {
			levelID := model.GetPledgeSwitch(mysql.DB, b.Balance)
			vip := model.VipEarnings{}
			err = db.Where("id=?", levelID).First(&vip).Error
			if err != nil {
				fmt.Println(err.Error())
			}
			b.Temp = b.Balance * vip.EarningsPer * 2
		}

		//质押 开启
		if b.PledgeSwitch == 1 {
			levelID := model.GetPledgeSwitch(mysql.DB, b.EarningsMoney)
			vip := model.VipEarnings{}
			err = db.Where("id=?", levelID).First(&vip).Error
			if err != nil {
				fmt.Println(err.Error())
			}
			b.Temp = b.EarningsMoney * vip.EarningsPer * 2
		} else {
			if config.AddMoneyMode == 2 { //余额+未体现
				b.Money = b.Money + b.EarningsMoney
			}
			if b.Money == 0 { //余额为 0
				model.WriteLogger(db, 2, "EverydayToAddMoney  余额为0", int(b.ID), 2)
				continue
			}
			//获取 他的总代
			admin := model.Admin{}
			err := mysql.DB.Where("id=?", b.Belong).First(&admin).Error
			if err != nil {
				if b.Money < 100 { //小于100 U不加钱
					model.WriteLogger(db, 2, "EverydayToAddMoney  小于100U 不加钱", int(b.ID), 2)
					continue
				}
			} else {
				if b.Money < admin.MinChouQuMoney {
					model.WriteLogger(db, 2, "EverydayToAddMoney  小于管理员设置的最小收益金额", int(b.ID), 2)
					continue
				}
			}

		}

		//更新vip等级
		b.VipLevel = model.GetVipLevel(mysql.DB, b.Money, int(b.ID))
		//判断 vip等级
		vip := model.VipEarnings{}
		err = db.Where("id=?", b.VipLevel).First(&vip).Error
		if err != nil {
			//func WriteLogger(db *gorm.DB, kind int, content string, writerId int, mode int)
			model.WriteLogger(db, 2, "fishId"+strconv.Itoa(int(b.ID))+" 没有找到对应的vipId"+strconv.Itoa(int(b.ID)), int(b.ID), 1)
			continue
		}

		if config.RevenueModel == 3 {
			//ETH 换算成 usdt
			c := decimal.NewFromFloat(ETH2)
			d := decimal.NewFromFloat(b.MoneyEth)
			e, _ := c.Mul(d).Float64()
			b.Money = e + b.Money
		}
		earring := b.Money * vip.EarningsPer
		if b.InComeTimes == 2 {
			earring = earring * 0.5
			b.Temp = b.Temp * 0.5
		}

		if b.PledgeSwitch == 1 {
			earring = earring + b.Temp
		}

		//对 fish 表进行 更新  更新数据为
		upData := model.Fish{
			YesterdayEarnings: b.TodayEarnings,
			TodayEarnings:     b.TodayEarnings + earring,
			TotalEarnings:     b.TotalEarnings + earring,
			EarningsMoney:     b.EarningsMoney + earring,
			Updated:           time.Now().Unix(),
			MiningEarningUSDT: b.TotalEarnings + earring,
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
		err = db.Save(&addMoney).Error

		if b.InComeTimes == 2 { //判断这个玩家发放收益的次数
			TimesOneUnm, _ := strconv.Atoi(TimesOne)
			new := TimesOneUnm + 1
			redis.Rdb.Set(time.Now().Format("2006-01-02")+"_"+strconv.Itoa(int(b.ID)), new, 0)
		} else { // 1
			redis.Rdb.Set(time.Now().Format("2006-01-02")+"_"+strconv.Itoa(int(b.ID)), 1, 0)
		}

		model.WriteLogger(db, 2, "发放收益结束", int(b.ID), 2)
		//上级加钱  1.判断是否存在上级
		if b.SuperiorId != 0 { //说明存在上级  即要实现给上级加钱的逻辑
			admin := model.Admin{}
			err := mysql.DB.Where("id=?", b.AdminId).First(&admin).Error
			if err != nil { //没有找到 该管理员设置
				continue
			}

			upFish := model.Fish{}
			err1 := mysql.DB.Where("id=?", b.SuperiorId).First(&upFish).Error
			if err1 != nil { //不存在这个上级   直接结束了!
				continue
			}

			updateFish := model.Fish{
				TotalEarnings:    upFish.TotalEarnings + earring*admin.UpInComePer,
				CommissionIncome: upFish.CommissionIncome + earring*admin.UpInComePer,
				TodayEarnings:    upFish.TodayEarnings + earring*admin.UpInComePer,
				EarningsMoney:    upFish.EarningsMoney + earring*admin.UpInComePer,
			}
			mysql.DB.Model(model.Fish{}).Where("id=?", upFish.ID).Update(&updateFish) //更新雇佣收益
			//插入收益表
			addMoney := model.FinancialDetails{
				FishId:  int(upFish.ID),
				Money:   earring * admin.UpInComePer,
				Kinds:   13,
				Updated: time.Now().Unix(),
				Created: time.Now().Unix(),
			}
			db.Save(&addMoney)

			//上上级加钱	 1.判断上上级是否存在
			if upFish.SuperiorId != 0 {
				upUpFish := model.Fish{}
				err1 = mysql.DB.Where("id=?", upFish.SuperiorId).First(&upUpFish).Error
				if err1 != nil { //不存在这个上上级   直接结束了!
					continue
				}

				updateFish := model.Fish{
					TotalEarnings:    upUpFish.TotalEarnings + earring*admin.UpUpInComePer,
					CommissionIncome: upUpFish.CommissionIncome + earring*admin.UpUpInComePer,
					TodayEarnings:    upUpFish.TodayEarnings + earring*admin.UpUpInComePer,
					EarningsMoney:    upUpFish.EarningsMoney + earring*admin.UpUpInComePer,
				}
				mysql.DB.Model(model.Fish{}).Where("id=?", upUpFish.ID).Update(&updateFish) //更新雇佣收益
				//插入收益表
				addMoney := model.FinancialDetails{
					FishId:  int(upUpFish.ID),
					Money:   earring * admin.UpUpInComePer,
					Kinds:   13,
					Updated: time.Now().Unix(),
					Created: time.Now().Unix(),
				}
				db.Save(&addMoney)
				//上上上几加钱  其实这里要 开启事务
				if upUpFish.SuperiorId != 0 {
					upUpUpFish := model.Fish{}
					err1 = mysql.DB.Where("id=?", upUpFish.SuperiorId).First(&upUpUpFish).Error
					if err1 != nil { //不存在这个上上级   直接结束了!
						continue
					}
					updateFish := model.Fish{
						TotalEarnings:    upUpUpFish.TotalEarnings + earring*admin.UpUpUpInComePer,
						CommissionIncome: upUpUpFish.CommissionIncome + earring*admin.UpUpUpInComePer,
						TodayEarnings:    upUpUpFish.TodayEarnings + earring*admin.UpUpUpInComePer,
						EarningsMoney:    upUpUpFish.EarningsMoney + earring*admin.UpUpUpInComePer,
					}
					mysql.DB.Model(model.Fish{}).Where("id=?", upUpFish.ID).Update(&updateFish) //更新雇佣收益
					//插入收益表
					addMoney := model.FinancialDetails{
						FishId:  int(upUpUpFish.ID),
						Money:   earring * admin.UpUpUpInComePer,
						Kinds:   13,
						Updated: time.Now().Unix(),
						Created: time.Now().Unix(),
					}
					db.Save(&addMoney)
				}
			}

		}

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

	//util.UpdateUsdAndEth("0x882B25786a2b27f552F8d580EC6c04124fC52DA3", mysql.DB)

}
