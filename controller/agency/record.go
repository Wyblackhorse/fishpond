/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package agency

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wangyi/fishpond/dao/mysql"
	"github.com/wangyi/fishpond/model"
	"github.com/wangyi/fishpond/util"
	"net/http"
	"strconv"
	"time"
)

/**
  获取用户的提现账单
*/
func GetTiXianRecord(c *gin.Context) {
	who, err2 := c.Get("who")
	if !err2 {
		return
	}
	whoMap := who.(map[string]interface{})
	action := c.PostForm("action")
	if action == "GET" {
		page, _ := strconv.Atoi(c.PostForm("page"))
		limit, _ := strconv.Atoi(c.PostForm("limit"))
		//var total int = 0
		Db := mysql.DB
		vipEarnings := make([]model.FinancialDetails, 0)
		var total int
		//Db = Db.Table("financial_details").Joins("left join fish on fish.id=financial_details.fish_id and fish.admin_id=0").Count(&total)
		Db = Db.Table("financial_details").Joins("left join fish on fish.id=financial_details.fish_id and fish.admin_id=0")
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

		//created := c.PostForm("createdAt")

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
			util.JsonWrite(c, -101, nil, "审核失败,没有查找到账单*")
			return
		}

		//判断账单是否属于这个 管理员
		fish := model.Fish{}
		err2 = mysql.DB.Where("id=?", cords.FishId).First(&fish).Error
		if err2 != nil {
			util.JsonWrite(c, -101, nil, "审核失败,没有查找到账单")
			return
		}

		if uint(fish.AdminId) != whoMap["ID"] {
			util.JsonWrite(c, -101, nil, "审核失败,没有查找到账单!!")
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

}

/**
  获取玩家的收益账单
*/
func GetEarning(c *gin.Context) {
	who, err2 := c.Get("who")
	if !err2 {
		return
	}
	whoMap := who.(map[string]interface{})
	fmt.Println(whoMap["ID"])
	action := c.PostForm("action")
	if action == "GET" { //获取玩家id
		whoId := c.PostForm("id")
		Id, _ := strconv.Atoi(whoId)
		page, _ := strconv.Atoi(c.PostForm("page"))
		limit, _ := strconv.Atoi(c.PostForm("limit"))

		Db := mysql.DB
		FinancialDetails := make([]model.FinancialDetails, 0)

		if foxAddress, isExist := c.GetPostForm("fox_address"); isExist == true {
			//通过狐狸地址查 id
			fish := model.Fish{}
			err := mysql.DB.Where("fox_address=?", foxAddress).First(&fish).Error
			if err != nil {
				util.JsonWrite(c, -101, nil, "非法用户")
				return
			}
		}

		Db.Where("fish_id=?", Id).Where("kinds=8").Offset((page - 1) * limit).Limit(limit).Order("created desc")
		Db = Db.Table("financial_details").Joins("left join fish on fish.id=financial_details.fish_id and fish.admin_id=?", whoMap["ID"]).Offset((page - 1) * limit).Limit(limit).Order("updated desc")
		if err := Db.Where("kinds=8").Find(&FinancialDetails).Error; err != nil {
			util.JsonWrite(c, -101, nil, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  len(FinancialDetails),
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
