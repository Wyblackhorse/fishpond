/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package agency

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi/fishpond/dao/redis"
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
