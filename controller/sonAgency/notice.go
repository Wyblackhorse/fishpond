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
	"github.com/wangyi/fishpond/model"
	"github.com/wangyi/fishpond/util"
)

func SetNotice(c *gin.Context) {

	who, _ := c.Get("who")

	whoMap := who.(map[string]interface{})

	Notice := c.Query("Notice")

	admin := model.Admin{}
	err := mysql.DB.Where("id=?", whoMap["ID"]).First(&admin).Error
	if err != nil {
		util.JsonWrite(c, -101, nil, "非法请求")
		return
	}

	mysql.DB.Model(&model.Admin{}).Where("id=?", admin.ID).Update("Notice", Notice)

	util.JsonWrite(c, 200, nil, "修改成功")

	return

}
