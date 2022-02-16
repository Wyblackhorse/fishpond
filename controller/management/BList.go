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
)

func GetBList(c *gin.Context) {

	action := c.PostForm("action")
	if action == "GET" {
		data := make([]model.BAddressList, 0)
		err := mysql.DB.Find(&data).Error
		if err != nil {
			util.JsonWrite(c, -101, nil, "获取失败")
			return
		}
		util.JsonWrite(c, 200, data, "获取成功")
		return
	}
}



