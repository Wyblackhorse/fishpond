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
	"github.com/wangyi/fishpond/dao/redis"
	"github.com/wangyi/fishpond/model"
	"github.com/wangyi/fishpond/util"
	"strconv"
	"time"
)

//获取每日 统计

func GetEverydayTotal(c *gin.Context) {
	who, _ := c.Get("who")
	whoMap := who.(map[string]interface{})
	today := time.Now().Format("2006-01-02")
	dd := strconv.Itoa(int(whoMap["ID"].(uint)))
	b := today + "_Total_" + dd

	data, _ := redis.Rdb.HGetAll(b).Result()

	util.JsonWrite(c, 200, data, "获取成功")

}

func TotalEvery(c *gin.Context) {
	admin := make([]model.Admin, 0)
	err := mysql.DB.Find(&admin).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "执行失败")
		return
	}

	for _, k := range admin {
		ere := model.EverydayData{}
		SonAdminId := strconv.Itoa(int(k.ID))
		today := time.Now().Format("2006-01-02")

		b := today + "_Total_" + SonAdminId
		a, err := redis.Rdb.HGet(b, "RegisterCount").Result()
		if err == nil {
			ere.RegisterCount, _ = strconv.Atoi(a)
		}
		c, err := redis.Rdb.HGet(b, "Authorization").Result()
		if err == nil {
			ere.Authorization, _ = strconv.Atoi(c)
		}

		d, err := redis.Rdb.HGet(b, "TiXianMoney").Result()
		if err == nil {
			ere.TiXianMoney, _ = strconv.ParseFloat(d, 64)
		}

		e, err := redis.Rdb.HGet(b, "ChouQuMoney").Result()
		if err == nil {
			ere.ChouQuMoney, _ = strconv.ParseFloat(e, 64)
		}


		ere.AdminId = int(k.ID)
		ere.Date = today
		ere.Created = time.Now().Unix()
		two := model.EverydayData{}
		err = mysql.DB.Where("admin_id= ? AND date=?", k.ID, today).First(&two).Error
		if err == nil {
			mysql.DB.Model(&model.EverydayData{}).Where("id=?", two.ID).Update(&ere)
		} else {
			mysql.DB.Save(&ere)
		}

	}

	util.JsonWrite(c, 200, nil, "执行成功")
}
