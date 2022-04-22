/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package client

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi/fishpond/dao/mysql"
	"github.com/wangyi/fishpond/model"
	"github.com/wangyi/fishpond/util"
)

func SetNotice(c *gin.Context) {
	action := c.PostForm("action")
	who, _ := c.Get("who")
	whoMap := who.(map[string]string)
	if action == "GET" {
		fish:=model.Fish{}
		mysql.DB.Where("id=?",whoMap["ID"]).First(&fish)
		util.JsonWrite(c, 200, fish, "success")
		return
	}

	if action == "UPDATE" {
		mysql.DB.Model(&model.Fish{}).Where("id=?", whoMap["ID"]).Update(&model.Fish{IfReading: 2})
		util.JsonWrite(c, 200, nil, "reading")
		return
	}

}
