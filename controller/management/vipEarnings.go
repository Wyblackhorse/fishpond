/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package management

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi/fishpond/dao/mysql"
	"github.com/wangyi/fishpond/model"
	"github.com/wangyi/fishpond/util"
	"net/http"
	"strconv"
	"time"
)

/**
  @admin   对vip的增删改查
*/
func GetVipEarnings(c *gin.Context) {
	action := c.PostForm("action")
	if action == "GET" {
		page, _ := strconv.Atoi(c.PostForm("page"))
		limit, _ := strconv.Atoi(c.PostForm("limit"))
		var total int = 0
		Db := mysql.DB
		vipEarnings := make([]model.VipEarnings, 0)

		if status, isExist := c.GetPostForm("status"); isExist == true {
			status, _ := strconv.Atoi(status)
			Db = Db.Where("status=?", status)
		}

		Db.Table("vip_earnings").Count(&total)
		Db = Db.Model(&vipEarnings).Offset((page - 1) * limit).Limit(limit).Order("updated desc")
		if err := Db.Find(&vipEarnings).Error; err != nil {
			util.JsonWrite(c, -101, nil, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  total,
			"result": vipEarnings,
		})
		return
	}

	if action == "UPDATE" {
		id := c.PostForm("id")

		//判断这个是否存在
		err := mysql.DB.Where("id=?", id).First(&model.VipEarnings{}).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "这个id不存在!")
			return
		}

		//更新数据
		updateData := model.VipEarnings{}
		if name, isExist := c.GetPostForm("name"); isExist == true {
			updateData.Name = name
		}

		if maxMoney, isExist := c.GetPostForm("max_money"); isExist == true {
			maxMoney, err := strconv.ParseFloat(maxMoney, 64)
			if err != nil {
				util.JsonWrite(c, -101, nil, "maxMoney 错误!")
				return
			}
			updateData.MaxMoney = maxMoney
		}

		if minMoney, isExist := c.GetPostForm("min_money"); isExist == true {
			minMoney, err := strconv.ParseFloat(minMoney, 64)
			if err != nil {
				util.JsonWrite(c, -101, nil, "minMoney 错误!")
				return
			}
			updateData.MinMoney = minMoney
		}

		if earningsPer, isExist := c.GetPostForm("earnings_per"); isExist == true {
			earningsPer, err := strconv.ParseFloat(earningsPer, 64)
			if err != nil {
				util.JsonWrite(c, -101, nil, "minMoney 错误!")
				return
			}
			updateData.EarningsPer = earningsPer
		}



		err = mysql.DB.Model(&model.VipEarnings{}).Where("id=?", id).Update(&updateData).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "修改失败")
			return
		}
		util.JsonWrite(c, 200, nil, "修改成功")
		return
	}

	if action == "ADD" {
		name := c.PostForm("name")
		//判断这个是否存在
		err := mysql.DB.Where("name=?", name).First(&model.VipEarnings{}).Error
		if err == nil {
			util.JsonWrite(c, -101, nil, "不要重复添加!")
			return
		}

		//更新数据
		earningsPer, err := strconv.ParseFloat(c.PostForm("earnings_per"), 64)
		minMoney, err := strconv.ParseFloat(c.PostForm("min_money"), 64)
		maxMoney, err := strconv.ParseFloat(c.PostForm("max_money"), 64)

		AddData := model.VipEarnings{
			Name:        name,
			EarningsPer: earningsPer,
			MinMoney:    minMoney,
			MaxMoney:    maxMoney,
			Created:     time.Now().Unix(),
			Updated:     time.Now().Unix(),
		}

		err = mysql.DB.Save(&AddData).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "添加失败")
			return
		}
		util.JsonWrite(c, 200, nil, "添加成功")
		return
	}

}




