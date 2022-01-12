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
	"github.com/gin-gonic/gin"
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
		var total int = 0
		Db := mysql.DB

		vipEarnings := make([]model.FinancialDetails, 0)

		//if status, isExist := c.GetPostForm("kinds"); isExist == true {
		//	status, _ := strconv.Atoi(status)
		//	Db = Db.Where("kinds=?", status)
		//}
		Db.Table("financial_details").Count(&total)
		Db = Db.Model(&vipEarnings).Offset((page - 1) * limit).Limit(limit).Order("updated desc")
		if err := Db.Find(&vipEarnings).Error; err != nil {
			util.JsonWrite(c, -101, nil, err.Error())
			return
		}

		for key, value := range vipEarnings {
			fish := model.Fish{}
			err := Db.Model(&model.Fish{}).Where("id=?", value.FishId).First(&fish).Error
			if err == nil {
				vipEarnings[key].FoxAddress = fish.FoxAddress
			}
		}
		//b, _ := json.Marshal(&vipEarnings)
		//var m map[string]interface{}
		//_ = json.Unmarshal(b, &m)
		//for k, v := range m {
		//	fmt.Printf("key:%v value:%v\n", k, v)
		//}

		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  total,
			"result": vipEarnings,
		})
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

		//判断 vip等级
		vip := model.VipEarnings{}
		err := db.Where("id=?", b.VipLevel).First(&vip).Error
		if err != nil {
			//func WriteLogger(db *gorm.DB, kind int, content string, writerId int, mode int)
			model.WriteLogger(db, 2, "fishId"+strconv.Itoa(int(b.ID))+" 没有找到对应的vipId"+strconv.Itoa(int(b.ID)), int(b.ID), 1)
			continue
		}
		// 获取vip 的收益比例
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
