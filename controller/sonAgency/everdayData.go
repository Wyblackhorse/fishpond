/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package sonAgency

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

/**
  统计数据分析
*/
func GetTotal(c *gin.Context) {

	who, _ := c.Get("who")
	MapWho := who.(map[string]interface{})

	action := c.PostForm("action")
	if action == "GET" {
		page, _ := strconv.Atoi(c.PostForm("page"))
		limit, _ := strconv.Atoi(c.PostForm("limit"))
		var total int = 0
		ere := make([]model.EverydayData, 0)

		Db := mysql.DB.Where("admin_id=?", MapWho["ID"])

		if date, isE := c.GetPostForm("date"); isE == true {
			Db = Db.Where("date=?", date)
		}

		Db.Model(&model.EverydayData{}).Count(&total)
		Db = Db.Model(&ere).Offset((page - 1) * limit).Limit(limit).Order("created desc").Find(&ere)
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  total,
			"result": ere,
		})
		return

	}
}
