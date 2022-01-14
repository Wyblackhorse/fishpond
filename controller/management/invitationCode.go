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
	"github.com/wangyi/fishpond/dao/redis"
	"github.com/wangyi/fishpond/util"
	"time"
)

/****
  邀请码
*/
type InvitationCode struct {
	InvitationCode string
	Token          string
}

type AllInvitationCode struct {
	Data InvitationCode
}

func SetInvitationCode(c *gin.Context) {
	action := c.PostForm("action")

	//获取邀请码
	if action == "GET" {
		data, _ := redis.Rdb.HGetAll("InvitationCode").Result()
		var mySliceMap []map[string]string

		for k, v := range data {
			countryCapitalMap := make(map[string]string)
			countryCapitalMap["key"] = k
			countryCapitalMap["value"] = v
			mySliceMap = append(mySliceMap, countryCapitalMap)

		}

		util.JsonWrite(c, 200, mySliceMap, "获取成功")
		return
	}

	if action == "ADD" {
		s := "A" + time.Now().Format("2006-01-02") + util.RandStr(20)
		redis.Rdb.HSet("InvitationCode", s, c.PostForm("token"))
		util.JsonWrite(c, 200, nil, "生成成功")
		return
	}

}
